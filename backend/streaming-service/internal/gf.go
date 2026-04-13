package main

import (
	"context"
	"fmt"
	"io"
	"music/streaming-service/proto/streaming"
	"music/track-service/proto/track"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type StreamingServer struct {
	streaming.UnimplementedStreamingServiceServer
	trackClient track.TrackServiceClient
}

func NewStreamingServer(ctx context.Context, trackServiceAddr string) (*StreamingServer, error) {

	conn, err := grpc.NewClient(trackServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to track service: %w", err)
	}

	return &StreamingServer{
		trackClient: track.NewTrackServiceClient(conn),
	}, nil
}

func (s *StreamingServer) StreamTrack(req *streaming.StreamRequest, stream streaming.StreamingService_StreamTrackServer) error {
	ctx := stream.Context()

	track, err := s.trackClient.GetTrack(ctx, &track.GetTrackRequest{
		Id: req.Id,
	})
	if err != nil {
		return status.Errorf(codes.NotFound, "track not found: %v", err)
	}

	file, err := os.Open(track.Track.AudioUrl)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to open file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get file info: %v", err)
	}
	fileSize := stat.Size()

	if req.StartPosition > 0 {
		_, err = file.Seek(int64(req.StartPosition), io.SeekStart)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid start position: %v", err)
		}
	}

	mimeType := getMimeTypeFromFile(track.Track.AudioUrl)

	chunkSize := 64 * 1024 // 64KB по умолчанию
	if req.ChunkSize > 0 && req.ChunkSize <= 1024*1024 {
		chunkSize = int(req.ChunkSize)
	}

	buffer := make([]byte, chunkSize)
	position := req.StartPosition

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := file.Read(buffer)

		// Файл закончился
		if err == io.EOF {
			if n > 0 {
				// Отправляем последние данные
				chunk := &streaming.Chunk{
					Data:      buffer[:n],
					Position:  position,
					IsLast:    true,
					MimeType:  mimeType,
					TotalSize: fileSize,
				}
				if err := stream.Send(chunk); err != nil {
					return err
				}
			}
			break // Выходим из цикла
		}

		// Реальная ошибка
		if err != nil {
			return status.Errorf(codes.Internal, "error reading file: %v", err)
		}

		// Отправляем обычный чанк
		chunk := &streaming.Chunk{
			Data:      buffer[:n],
			Position:  position,
			IsLast:    false,
			MimeType:  mimeType,
			TotalSize: fileSize,
		}

		if err := stream.Send(chunk); err != nil {
			return err
		}

		position += int64(n)
	}

	return nil
}

func getMimeTypeFromFile(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".flac":
		return "audio/flac"
	case ".aac":
		return "audio/aac"
	case ".m4a":
		return "audio/mp4"
	case ".opus":
		return "audio/opus"
	default:
		// Пробуем определить по содержимому
		file, err := os.Open(filePath)
		if err != nil {
			return "application/octet-stream"
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			return "application/octet-stream"
		}

		return http.DetectContentType(buffer)
	}
}
