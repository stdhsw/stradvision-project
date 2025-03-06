package producer

import (
	"time"

	"github.com/IBM/sarama"
)

type producerConfig struct {
	config *sarama.Config

	errFunc     func(ts time.Time, topic string, partition int32, err error)
	successFunc func(ts time.Time, topic string, partition int32)
}

func defaultConfig() *producerConfig {
	config := sarama.NewConfig()
	// 기본 설정 (사용자 설정 가능)
	config.Producer.MaxMessageBytes = 1024 * 1024             // 최대 메시지 크기 1MB
	config.Producer.RequiredAcks = sarama.WaitForLocal        // 리더 브로커만 ACK 반환
	config.Producer.Timeout = 3 * time.Second                 // 3초 타임아웃
	config.Producer.Retry.Max = 5                             // 최대 재시도 횟수 5회
	config.Producer.Retry.Backoff = 100 * time.Millisecond    // 재시도 간격 100ms
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 랜덤 파티션 선택

	// 메시지 압축 설정 (사용자 설정 가능)
	config.Producer.Compression = sarama.CompressionSnappy            // 메시지 압축
	config.Producer.CompressionLevel = sarama.CompressionLevelDefault // 압축 레벨

	// Flush 설정 (사용자 설정 가능)
	config.Producer.Flush.Frequency = 500 * time.Millisecond // 0.5초마다 전송
	config.Producer.Flush.Bytes = 1024 * 1024                // 1MB마다 전송
	config.Producer.Flush.MaxMessages = 1000                 // 1000개마다 전송

	// Transaction 설정 (사용자 설정 불가능)
	config.Producer.Transaction.Timeout = 5 * time.Second              // 트랜잭션 타임아웃 5초
	config.Producer.Transaction.Retry.Max = 5                          // 최대 재시도 횟수 5회
	config.Producer.Transaction.Retry.Backoff = 100 * time.Millisecond // 재시도 간격 100ms

	// ACK 반환 설정 (사용자 설정 불가능)
	config.Producer.Return.Successes = true // 성공 시 메시지 반환
	config.Producer.Return.Errors = true    // 전송 실패 시 에러 반환

	pConfig := &producerConfig{
		config:      config,
		errFunc:     func(ts time.Time, topic string, partition int32, err error) {},
		successFunc: func(ts time.Time, topic string, partition int32) {},
	}

	return pConfig
}

// fromOptions producer 설정 반환
func fromOptions(options []Option) *producerConfig {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	return config
}

type Option func(*producerConfig)

// WithErrorFunc 에러 콜백 함수 설정
func WithErrorFunc(errFunc func(ts time.Time, topic string, partition int32, err error)) Option {
	return func(pConfig *producerConfig) {
		if errFunc != nil {
			pConfig.errFunc = errFunc
		}
	}
}

// WithSuccessFunc 성공 콜백 함수 설정
func WithSuccessFunc(successFunc func(ts time.Time, topic string, partition int32)) Option {
	return func(pConfig *producerConfig) {
		if successFunc != nil {
			pConfig.successFunc = successFunc
		}
	}
}

// WithMaxMessageBytes 메시지 최대 크기 설정
func WithMaxMessageBytes(maxMessageBytes int) Option {
	return func(pConfig *producerConfig) {
		if maxMessageBytes > 0 {
			pConfig.config.Producer.MaxMessageBytes = maxMessageBytes
		}
	}
}

// WithRequiredAcks ACK 반환 설정
// 0: NoResponse, 1: WaitForLocal, -1: WaitForAll
// default: WaitForLocal
func WithRequiredAcks(requiredAcks int16) Option {
	return func(pConfig *producerConfig) {
		ack := sarama.RequiredAcks(requiredAcks)
		switch ack {
		case sarama.WaitForAll, sarama.WaitForLocal, sarama.NoResponse:
			pConfig.config.Producer.RequiredAcks = ack
		default:
			pConfig.config.Producer.RequiredAcks = sarama.WaitForLocal
		}
	}
}

// WithTimeout 타임아웃 설정
// timeout: 최소값 1초
func WithTimeout(timeout time.Duration) Option {
	return func(pConfig *producerConfig) {
		if timeout > time.Second {
			pConfig.config.Producer.Timeout = timeout
		}
	}
}

// WithRetry 최대 재시도 횟수 설정
// max: 최소값 0
func WithRetry(max int) Option {
	return func(pConfig *producerConfig) {
		if max >= 0 {
			pConfig.config.Producer.Retry.Max = max
		}
	}
}

// WithRetryBackoff 재시도 간격 설정
// backoff: 최소값 0
func WithRetryBackoff(backoff time.Duration) Option {
	return func(pConfig *producerConfig) {
		if backoff >= 0 {
			pConfig.config.Producer.Retry.Backoff = backoff
		}
	}
}

// WithPartitioner 파티션 선택 설정
// partitioner: 0(Random), 1(RoundRobin), 2(Hash)
// default: Random
func WithPartitioner(partitioner int) Option {
	return func(pConfig *producerConfig) {
		switch partitioner {
		case 0:
			pConfig.config.Producer.Partitioner = sarama.NewRandomPartitioner
		case 1:
			pConfig.config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		case 2:
			pConfig.config.Producer.Partitioner = sarama.NewHashPartitioner
		default:
			pConfig.config.Producer.Partitioner = sarama.NewRandomPartitioner
		}
	}
}

// WithCompression 메시지 압축 설정
// compression: 0(None), 1(GZIP), 2(Snappy), 3(LZ4), 4(ZSTD)
// default: None
func WithCompression(compression int) Option {
	return func(pConfig *producerConfig) {
		switch compression {
		case 0:
			pConfig.config.Producer.Compression = sarama.CompressionNone
		case 1:
			pConfig.config.Producer.Compression = sarama.CompressionGZIP
		case 2:
			pConfig.config.Producer.Compression = sarama.CompressionSnappy
		case 3:
			pConfig.config.Producer.Compression = sarama.CompressionLZ4
		case 4:
			pConfig.config.Producer.Compression = sarama.CompressionZSTD
		default:
			pConfig.config.Producer.Compression = sarama.CompressionNone
		}
	}
}

// WithCompressionLevel 압축 레벨 설정
// level: -1(Default), 0~9
// default: Default
func WithCompressionLevel(level int) Option {
	return func(pConfig *producerConfig) {
		if level >= -1 {
			pConfig.config.Producer.CompressionLevel = int(level)
		}
	}
}

// WithFlushFrequency 전송 주기 설정
// frequency: 최소값 0
func WithFlushFrequency(frequency time.Duration) Option {
	return func(pConfig *producerConfig) {
		if frequency >= 0 {
			pConfig.config.Producer.Flush.Frequency = frequency
		}
	}
}

// WithFlushBytes 전송 크기 설정
// bytes: 최소값 0
func WithFlushBytes(bytes int) Option {
	return func(pConfig *producerConfig) {
		if bytes >= 0 {
			pConfig.config.Producer.Flush.Bytes = int(bytes)
		}
	}
}

// WithFlushMaxMessages 전송 개수 설정
// max: 최소값 0
func WithFlushMaxMessages(max int) Option {
	return func(pConfig *producerConfig) {
		if max >= 0 {
			pConfig.config.Producer.Flush.MaxMessages = max
		}
	}
}
