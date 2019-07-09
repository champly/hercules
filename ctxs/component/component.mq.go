package component

type IMQ interface {
	Produce(queueName, value string) error
}
