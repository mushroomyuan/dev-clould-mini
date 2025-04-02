package event_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// KafkaTestSuite 包含测试所需的配置和客户端
type KafkaTestSuite struct {
	brokers     []string
	adminClient *kafka.Client
	ctx         context.Context
	cancel      context.CancelFunc
}

// SetupSuite 设置测试环境
func (s *KafkaTestSuite) SetupSuite(t *testing.T) {
	// 从环境变量获取 broker 列表，如果没有则使用默认值
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092,localhost:9093,localhost:9094"
	}
	s.brokers = strings.Split(brokers, ",")

	// 创建上下文
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 10*time.Second)

	// 创建 AdminClient
	client, err := kafka.DialContext(s.ctx, "tcp", s.brokers[0])
	require.NoError(t, err, "无法连接到 Kafka")
	s.adminClient = client
}

// TearDownSuite 清理测试环境
func (s *KafkaTestSuite) TearDownSuite(t *testing.T) {
	if s.adminClient != nil {
		s.adminClient.Close()
	}
	if s.cancel != nil {
		s.cancel()
	}
}

// TestCreateTopic 测试创建主题
func TestCreateTopic(t *testing.T) {
	suite := &KafkaTestSuite{}
	suite.SetupSuite(t)
	defer suite.TearDownSuite(t)

	// 创建主题配置
	topicConfig := kafka.TopicConfig{
		Topic:             "maudit",
		NumPartitions:     3,
		ReplicationFactor: 3,
		ConfigEntries: []kafka.ConfigEntry{
			{
				ConfigName:  "cleanup.policy",
				ConfigValue: "delete",
			},
			{
				ConfigName:  "retention.ms",
				ConfigValue: "86400000", // 24小时
			},
		},
	}

	// 创建主题
	err := suite.adminClient.CreateTopics(suite.ctx, topicConfig)
	require.NoError(t, err, "创建主题失败")

	// 验证主题是否创建成功
	partitions, err := suite.adminClient.ReadPartitions()
	require.NoError(t, err, "获取分区信息失败")

	topics := make(map[string]int)
	for _, p := range partitions {
		topics[p.Topic]++
	}

	assert.Contains(t, topics, "maudit", "主题未创建成功")
	assert.Equal(t, 3, topics["maudit"], "分区数量不正确")
}

// TestWriteAndReadMessages 测试消息的写入和读取
func TestWriteAndReadMessages(t *testing.T) {
	suite := &KafkaTestSuite{}
	suite.SetupSuite(t)
	defer suite.TearDownSuite(t)

	// 创建写入器
	writer := &kafka.Writer{
		Addr:         kafka.TCP(suite.brokers...),
		Topic:        "maudit",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	defer writer.Close()

	// 写入消息
	messages := []kafka.Message{
		{
			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	}

	err := writer.WriteMessages(suite.ctx, messages...)
	require.NoError(t, err, "写入消息失败")

	// 创建读取器
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  suite.brokers,
		Topic:    "maudit",
		GroupID:  "test-group",
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  100 * time.Millisecond,
	})
	defer reader.Close()

	// 读取消息
	for i := 0; i < len(messages); i++ {
		msg, err := reader.ReadMessage(suite.ctx)
		require.NoError(t, err, "读取消息失败")
		assert.Equal(t, messages[i].Key, msg.Key, "消息 Key 不匹配")
		assert.Equal(t, messages[i].Value, msg.Value, "消息 Value 不匹配")
		fmt.Printf("读取到消息: topic=%s partition=%d offset=%d key=%s value=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}

// TestListTopics 测试获取主题列表
func TestListTopics(t *testing.T) {
	suite := &KafkaTestSuite{}
	suite.SetupSuite(t)
	defer suite.TearDownSuite(t)

	partitions, err := suite.adminClient.ReadPartitions()
	require.NoError(t, err, "获取分区信息失败")

	topics := make(map[string]int)
	for _, p := range partitions {
		topics[p.Topic]++
	}

	t.Log("当前主题列表:", topics)
}
