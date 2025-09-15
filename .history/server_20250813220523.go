package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"example-cloud-agent/proto"

	"google.golang.org/grpc"
)

type receiverServer struct {
	proto.UnimplementedReceiverServer
	clients map[string]proto.Receiver_ConnectServer
}

func (s *receiverServer) Connect(stream proto.Receiver_ConnectServer) error {
	var clientID string

	// 클라이언트로부터 메시지 수신 및 응답
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("클라이언트 연결 종료: %v", err)
			if clientID != "" {
				delete(s.clients, clientID)
			}
			return err
		}

		// 첫 번째 메시지에서 클라이언트 ID 저장
		if clientID == "" {
			clientID = msg.GetAgentId()
			s.clients[clientID] = stream
			log.Printf("새 클라이언트 연결: %s", clientID)

			// 환영 메시지 전송
			welcomeMsg := &proto.Message{
				AgentId:     "server",
				MessageType: "welcome",
				Content:     fmt.Sprintf("환영합니다, %s!", clientID),
				Timestamp:   time.Now().Unix(),
			}
			stream.Send(welcomeMsg)
		}

		log.Printf("메시지 수신: [%s] %s - %s", msg.GetMessageType(), msg.GetAgentId(), msg.GetContent())

		// 메시지 타입에 따른 처리
		switch msg.GetMessageType() {
		case "connect":
			log.Printf("클라이언트 %s 연결됨", msg.GetAgentId())
		case "ping":
			// ping에 대한 pong 응답
			pongMsg := &proto.Message{
				AgentId:     "server",
				MessageType: "pong",
				Content:     "pong",
				Timestamp:   time.Now().Unix(),
			}
			stream.Send(pongMsg)
		case "heartbeat":
			// 하트비트 확인 응답
			heartbeatAck := &proto.Message{
				AgentId:     "server",
				MessageType: "heartbeat_ack",
				Content:     "received",
				Timestamp:   time.Now().Unix(),
			}
			stream.Send(heartbeatAck)
		case "status_report":
			log.Printf("상태 보고 수신: %s", msg.GetContent())
		case "metric":
			log.Printf("메트릭 수신: %s", msg.GetContent())
		case "log":
			log.Printf("로그 수신: %s", msg.GetContent())
		case "event":
			log.Printf("이벤트 수신: %s", msg.GetContent())
		case "alert":
			log.Printf("알림 수신: %s", msg.GetContent())
		}

		// 주기적으로 명령 전송 (테스트용)
		if time.Now().Unix()%30 == 0 { // 30초마다
			commandMsg := &proto.Message{
				AgentId:     "server",
				MessageType: "command",
				Content:     fmt.Sprintf("테스트 명령 #%d", time.Now().Unix()),
				Timestamp:   time.Now().Unix(),
			}
			stream.Send(commandMsg)
		}
	}
}

func main() {
	fmt.Println("🚀 gRPC Receiver 서버 시작...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("포트 리스닝 실패: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterReceiverServer(s, &receiverServer{
		clients: make(map[string]proto.Receiver_ConnectServer),
	})

	fmt.Println("🔗 서버가 포트 50051에서 실행 중입니다...")
	fmt.Println("클라이언트 연결을 기다리는 중...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}
