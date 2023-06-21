package app

type Runnable interface {
	Run() error
}

type Option func(*application)
