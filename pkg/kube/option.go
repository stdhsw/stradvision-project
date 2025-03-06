package kube

import "time"

type clientConfig struct {
	kubeConfig string
	resyncTime time.Duration
}

func defaultConfig() *clientConfig {
	return &clientConfig{
		kubeConfig: "",
		resyncTime: 0,
	}
}

type Option func(*clientConfig)

func fromOptions(options []Option) *clientConfig {
	config := defaultConfig()
	for _, option := range options {
		option(config)
	}
	return config
}

// WithResyncTime resync 시간 설정
func WithResyncTime(resyncTime time.Duration) Option {
	return func(c *clientConfig) {
		if resyncTime > 0 {
			c.resyncTime = resyncTime
		}
	}
}

// WithKubeConfig kubeconfig 설정
func WithKubeConfig(kubeConfig string) Option {
	return func(c *clientConfig) {
		if kubeConfig != "" {
			c.kubeConfig = kubeConfig
		}
	}
}
