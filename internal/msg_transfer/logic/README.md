

1.Kafka 消费者组中可以存在多个消费者，Kafka 会以 partition 为单位将消息分给各个消费者。每条消息只会被消费者组的一个消费者消费。


online_history_msg 
(1)消费消息指令，将每1000条消息放入分发chan
(2)消息聚合指令，根据sendid或receiveid将消息聚合,在用hash的方式放入不同的chan
(3)将聚合后的消息在放入推送的mq