package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"example-receiver/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("☁️ Cloud Agent 클라이언트 시작...")
	fmt.Println("📋 모든 proto 메소드와 구조체를 사용하는 예제")

	// gRPC 서버에 연결
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("서버 연결 실패: %v", err)
	}
	defer conn.Close()

	// 클라이언트 생성
	client := proto.NewReceiverClient(conn)

	// 스트림 연결
	stream, err := client.Connect(context.Background())
	if err != nil {
		log.Fatalf("스트림 연결 실패: %v", err)
	}

	fmt.Println("🔗 Receiver 서버에 연결됨!")

	// 메시지 수신 루프
	go receiveMessages(stream)

	// 다양한 메시지 타입 전송 루프
	sendVariousMessages(stream)
}

func receiveMessages(stream proto.Receiver_ConnectClient) {
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("메시지 수신 오류: %v", err)
			return
		}

		fmt.Printf("📨 메시지 수신: [%s] %s - %s\n", msg.GetMessageType(), msg.GetAgentId(), msg.GetContent())

		// 모든 메시지 타입에 대한 처리
		switch msg.GetMessageType() {
		case "welcome":
			fmt.Println("🎉 서버로부터 환영 메시지 수신!")
		case "command":
			fmt.Printf("⚡ 명령 수신: %s\n", msg.GetContent())
			// 명령 실행 후 응답
			response := createMessage("cloud-agent-client", "command_response",
				fmt.Sprintf("명령 실행 완료: %s", msg.GetContent()))
			stream.Send(response)
		case "status_request":
			fmt.Println("📊 상태 요청 수신")
			statusResponse := createMessage("cloud-agent-client", "status_response",
				"CPU: 25%, Memory: 60%, Disk: 45%")
			stream.Send(statusResponse)
		case "config_update":
			fmt.Printf("⚙️ 설정 업데이트 수신: %s\n", msg.GetContent())
			configResponse := createMessage("cloud-agent-client", "config_ack",
				"설정 업데이트 완료")
			stream.Send(configResponse)
		case "heartbeat":
			fmt.Println("💓 하트비트 수신")
			heartbeatResponse := createMessage("cloud-agent-client", "heartbeat_ack", "alive")
			stream.Send(heartbeatResponse)
		}
	}
}

func sendVariousMessages(stream proto.Receiver_ConnectClient) {
	agentID := "cloud-agent-client"

	// 1. 연결 메시지 전송
	connectMsg := createMessage(agentID, "connect", "Cloud Agent 클라이언트 연결")
	if err := stream.Send(connectMsg); err != nil {
		log.Printf("연결 메시지 전송 실패: %v", err)
		return
	}
	fmt.Println("🔌 연결 메시지 전송")

	// 2. 초기화 메시지 전송
	initMsg := createMessage(agentID, "init", "클라이언트 초기화 완료")
	if err := stream.Send(initMsg); err != nil {
		log.Printf("초기화 메시지 전송 실패: %v", err)
		return
	}
	fmt.Println("🚀 초기화 메시지 전송")

	// 3. 등록 메시지 전송
	registerMsg := createMessage(agentID, "register", "agent_type=cloud,version=1.0.0")
	if err := stream.Send(registerMsg); err != nil {
		log.Printf("등록 메시지 전송 실패: %v", err)
		return
	}
	fmt.Println("📝 등록 메시지 전송")

	// 4. 설정 요청 메시지 전송
	configRequestMsg := createMessage(agentID, "config_request", "설정 정보 요청")
	if err := stream.Send(configRequestMsg); err != nil {
		log.Printf("설정 요청 메시지 전송 실패: %v", err)
		return
	}
	fmt.Println("⚙️ 설정 요청 메시지 전송")

	// 5. 상태 보고 메시지 전송
	statusMsg := createMessage(agentID, "status_report", "정상 동작 중")
	if err := stream.Send(statusMsg); err != nil {
		log.Printf("상태 보고 메시지 전송 실패: %v", err)
		return
	}
	fmt.Println("📊 상태 보고 메시지 전송")

	// 주기적으로 다양한 메시지 전송
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	messageCounter := 0
	for {
		select {
		case <-ticker.C:
			messageCounter++

			// 다양한 메시지 타입을 순환하며 전송
			switch messageCounter % 6 {
			case 0:
				// ping 메시지
				pingMsg := createMessage(agentID, "ping", "ping")
				if err := stream.Send(pingMsg); err != nil {
					log.Printf("ping 메시지 전송 실패: %v", err)
					return
				}
				fmt.Println("🏓 ping 메시지 전송")

			case 1:
				// 하트비트 메시지
				heartbeatMsg := createMessage(agentID, "heartbeat", "heartbeat")
				if err := stream.Send(heartbeatMsg); err != nil {
					log.Printf("하트비트 메시지 전송 실패: %v", err)
					return
				}
				fmt.Println("💓 하트비트 메시지 전송")

			case 2:
				// 로그 메시지
				logMsg := createMessage(agentID, "log", fmt.Sprintf("로그 메시지 #%d", messageCounter))
				if err := stream.Send(logMsg); err != nil {
					log.Printf("로그 메시지 전송 실패: %v", err)
					return
				}
				fmt.Println("📝 로그 메시지 전송")

			case 3:
				// 메트릭 메시지
				metricMsg := createMessage(agentID, "metric",
					fmt.Sprintf("cpu_usage=%d,memory_usage=%d,disk_usage=%d",
						messageCounter%100, (messageCounter*2)%100, (messageCounter*3)%100))
				if err := stream.Send(metricMsg); err != nil {
					log.Printf("메트릭 메시지 전송 실패: %v", err)
					return
				}
				fmt.Println("📈 메트릭 메시지 전송")

			case 4:
				// 이벤트 메시지
				eventMsg := createMessage(agentID, "event", fmt.Sprintf("이벤트 발생 #%d", messageCounter))
				if err := stream.Send(eventMsg); err != nil {
					log.Printf("이벤트 메시지 전송 실패: %v", err)
					return
				}
				fmt.Println("🎯 이벤트 메시지 전송")

			case 5:
				// 알림 메시지
				alertMsg := createMessage(agentID, "alert", fmt.Sprintf("알림 #%d", messageCounter))
				if err := stream.Send(alertMsg); err != nil {
					log.Printf("알림 메시지 전송 실패: %v", err)
					return
				}
				fmt.Println("🚨 알림 메시지 전송")
			}
		}
	}
}

// Message 구조체의 모든 메소드를 사용하는 헬퍼 함수
func createMessage(agentID, messageType, content string) *proto.Message {
	msg := &proto.Message{
		AgentId:     agentID,
		MessageType: messageType,
		Content:     content,
		Timestamp:   time.Now().Unix(),
	}

	// Message 구조체의 모든 getter 메소드 사용
	fmt.Printf("🔍 Message 구조체 메소드 테스트:\n")
	fmt.Printf("  - GetAgentId(): %s\n", msg.GetAgentId())
	fmt.Printf("  - GetMessageType(): %s\n", msg.GetMessageType())
	fmt.Printf("  - GetContent(): %s\n", msg.GetContent())
	fmt.Printf("  - GetTimestamp(): %d\n", msg.GetTimestamp())

	return msg
}

// Command 구조체의 모든 메소드를 사용하는 함수
func createCommand(commandID, commandType, targetAgentID, parameters string) *proto.Command {
	cmd := &proto.Command{
		CommandId:     commandID,
		CommandType:   commandType,
		TargetAgentId: targetAgentID,
		Parameters:    parameters,
	}

	// Command 구조체의 모든 getter 메소드 사용
	fmt.Printf("🔍 Command 구조체 메소드 테스트:\n")
	fmt.Printf("  - GetCommandId(): %s\n", cmd.GetCommandId())
	fmt.Printf("  - GetCommandType(): %s\n", cmd.GetCommandType())
	fmt.Printf("  - GetTargetAgentId(): %s\n", cmd.GetTargetAgentId())
	fmt.Printf("  - GetParameters(): %s\n", cmd.GetParameters())

	return cmd
}
