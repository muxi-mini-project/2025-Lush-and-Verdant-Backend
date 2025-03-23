package service

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"2025-Lush-and-Verdant-Backend/config"
	"2025-Lush-and-Verdant-Backend/dao"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type LikeService interface {
	SendMsg(like *request.ForestLikeReq) error
	GetLikes(to string) (string, error)
	GetForestLikeStatus(like *request.ForestLikeReq) bool
}

type LikeServiceImpl struct {
	Producer sarama.AsyncProducer
	Consumer sarama.Consumer
	jwt      *config.KafkaConfig
	Dao      dao.LikeDAO
}

func NewLikeServiceImpl(jwt *config.KafkaConfig, Dao dao.LikeDAO) *LikeServiceImpl {
	// 生产者
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	producerConfig.Producer.Compression = sarama.CompressionGZIP      //压缩
	producerConfig.Producer.Return.Successes = true                   // 回归
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll          //等待所有都有消息

	//配置生产者
	fmt.Print(jwt.Addr)
	producer, err := sarama.NewAsyncProducer([]string{jwt.Addr}, producerConfig)
	if err != nil {
		panic(err)
	}

	//配置消费者
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{jwt.Addr}, consumerConfig)
	if err != nil {
		panic(err)
	}
	return &LikeServiceImpl{
		Producer: producer,
		Consumer: consumer,
		jwt:      jwt,
		Dao:      Dao,
	}
}

// SendMsg 异步发送消息
func (lsr *LikeServiceImpl) SendMsg(like *request.ForestLikeReq) error {
	likeJson, err := json.Marshal(like)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "like_event",
		Key:   sarama.StringEncoder("forest"),
		Value: sarama.StringEncoder(likeJson),
	}
	lsr.Producer.Input() <- msg
	return nil
}

func (lsr *LikeServiceImpl) StartConsume() {
	partitions, err := lsr.Consumer.Partitions("like_event")
	if err != nil {
		log.Println("Error getting list of partitions: ", err)
	}
	for _, partition := range partitions {
		go func(partition int32) {
			partitionConsumer, err := lsr.Consumer.ConsumePartition("like_event", partition, sarama.OffsetNewest)
			if err != nil {
				log.Println(err)
			}
			defer partitionConsumer.Close()

			for msg := range partitionConsumer.Messages() {
				log.Printf("Received message: partition=%d, offset=%d, key=%s, value=%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//解析消息
				var like request.ForestLikeReq
				if err := json.Unmarshal(msg.Value, &like); err != nil {
					log.Println()
				}
				//处理消息
				lsr.HandLikeEvent(&like)
			}
		}(partition)

	}
}

func (lsr *LikeServiceImpl) HandLikeEvent(like *request.ForestLikeReq) {
	switch like.Action {
	case "like":
		lsr.ForestLike(like)
	case "unlike":
		lsr.ForestUnLike(like)
	default:
		log.Printf("未知的点赞消息 %+v\n", like)
	}
}
func (lsr *LikeServiceImpl) ForestLike(like *request.ForestLikeReq) {
	//先检查一下
	ok := lsr.Dao.Check(like)
	if !ok {
		//增加点赞数
		lsr.Dao.IncrementLike(like.To)
		lsr.Dao.SaveLike(like)
	} else {
		log.Printf("用户%s已经点赞用户%s的森林", like.From, like.To)
	}
}

func (lsr *LikeServiceImpl) ForestUnLike(like *request.ForestLikeReq) {
	// 检查一下
	ok := lsr.Dao.Check(like)
	if ok {
		lsr.Dao.DecrementLike(like.To)
		lsr.Dao.SaveUnlike(like)
	} else {
		log.Printf("用户%s还未点赞用户%s的森林", like.From, like.To)
	}
}

func (lsr *LikeServiceImpl) GetLikes(to string) (string, error) {
	return lsr.Dao.GetLikes(to)
}

func (lsr *LikeServiceImpl) GetForestLikeStatus(like *request.ForestLikeReq) bool {
	return lsr.Dao.Check(like)
}
