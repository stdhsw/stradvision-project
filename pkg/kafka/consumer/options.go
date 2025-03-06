package consumer

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type consumerConfig struct {
	config *sarama.Config

	doFunc  func([]byte)
	errFunc func(topic, msg string)
}

func defaultConfig() *consumerConfig {
	config := sarama.NewConfig()
	// 기본 설정 (사용자 설정 가능)
	config.Consumer.Fetch.Min = 1                              // 최소 메시지 크기 1바이트
	config.Consumer.Fetch.Default = 1024 * 1024                // 기본 메시지 크기 1MB
	config.Consumer.Retry.Backoff = 2 * time.Second            // 재시도 간격 2초
	config.Consumer.MaxWaitTime = 500 * time.Millisecond       // 최대 대기 시간 0.5초
	config.Consumer.MaxProcessingTime = 100 * time.Millisecond // 최대 처리 시간 0.1초

	// 그룹 설정 (사용자 설정 가능)
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky() // 리밸런스 전략
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second                   // 리밸런스 타임아웃
	config.Consumer.Group.Rebalance.Retry.Max = 4                                // 최대 재시도 횟수
	config.Consumer.Group.Rebalance.Retry.Backoff = 2 * time.Second              // 재시도 간격
	config.Consumer.Group.Session.Timeout = 10 * time.Second                     // 세션 타임아웃 10초
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second                   // 하트비트 간격 3초

	// 오프셋 설정 (사용자 설정 불가능)
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 가장 오래된 오프셋부터 시작
	config.Consumer.Offsets.AutoCommit.Enable = false     // 자동 커밋 비활성화

	cConfig := &consumerConfig{
		config:  config,
		doFunc:  func([]byte) {},
		errFunc: func(topic, msg string) {},
	}

	return cConfig
}

func fromOptions(options []Option) *consumerConfig {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	return config
}

type Option func(*consumerConfig)

// WithDoFunc 메시지 처리 함수 설정
func WithDoFunc(doFunc func([]byte)) Option {
	return func(c *consumerConfig) {
		if doFunc != nil {
			c.doFunc = doFunc
		}
	}
}

// WithErrFunc 에러 처리 함수 설정
func WithErrFunc(errFunc func(topic, msg string)) Option {
	return func(c *consumerConfig) {
		if errFunc != nil {
			c.errFunc = errFunc
		}
	}
}

// WithMinBytes 최소 메시지 크기 설정
func WithMinBytes(min int32) Option {
	return func(c *consumerConfig) {
		if min > 1 {
			c.config.Consumer.Fetch.Min = min
		}
	}
}

// WithMaxBytes 최대 메시지 크기 설정
func WithMaxBytes(max int32) Option {
	return func(c *consumerConfig) {
		if max > 1 {
			c.config.Consumer.Fetch.Default = max
		}
	}
}

// WithRetryBackoff 재시도 간격 설정
func WithRetryBackoff(backoff time.Duration) Option {
	return func(c *consumerConfig) {
		if backoff > 0 {
			c.config.Consumer.Retry.Backoff = backoff
		}
	}
}

// WithMaxWaitTime 최대 대기 시간 설정
func WithMaxWaitTime(wait time.Duration) Option {
	return func(c *consumerConfig) {
		if wait > 0 {
			c.config.Consumer.MaxWaitTime = wait
		}
	}
}

// WithMaxProcessingTime 최대 처리 시간 설정
func WithMaxProcessingTime(process time.Duration) Option {
	return func(c *consumerConfig) {
		if process > 0 {
			c.config.Consumer.MaxProcessingTime = process
		}
	}
}

// WithSessionTimeout 세션 타임아웃 설정
func WithSessionTimeout(timeout time.Duration) Option {
	return func(c *consumerConfig) {
		if timeout > 0 {
			c.config.Consumer.Group.Session.Timeout = timeout
		}
	}
}

// WithHeartbeatInterval 하트비트 간격 설정
func WithHeartbeatInterval(interval time.Duration) Option {
	return func(c *consumerConfig) {
		if interval > 0 {
			c.config.Consumer.Group.Heartbeat.Interval = interval
		}
	}
}

// WithRebalanceTimeout 리밸런스 타임아웃 설정
func WithRebalanceTimeout(timeout time.Duration) Option {
	return func(c *consumerConfig) {
		if timeout > 0 {
			c.config.Consumer.Group.Rebalance.Timeout = timeout
		}
	}
}

// WithRebalanceRetryMax 최대 재시도 횟수 설정
func WithRebalanceRetryMax(max int) Option {
	return func(c *consumerConfig) {
		if max >= 0 {
			c.config.Consumer.Group.Rebalance.Retry.Max = max
		}
	}
}

// WithRebalanceRetryBackoff 재시도 간격 설정
func WithRebalanceRetryBackoff(backoff time.Duration) Option {
	return func(c *consumerConfig) {
		if backoff > 0 {
			c.config.Consumer.Group.Rebalance.Retry.Backoff = backoff
		}
	}
}

// WithBalanceStrategy 리밸런스 전략 설정
// STICKY : 기본값, 파티션을 고정적으로 할당
// ROUNDROBIN : 파티션을 순차적으로 할당
// RANGE : 파티션을 범위에 따라 할당
func WithBalanceStrategy(strategy string) Option {
	return func(c *consumerConfig) {
		switch strings.ToUpper(strategy) {
		case "ROUNDROBIN":
			c.config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
		case "RANGE":
			c.config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
		case "STICKY":
			fallthrough
		default:
			c.config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
		}
	}
}
