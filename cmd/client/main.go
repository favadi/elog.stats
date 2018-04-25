package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "elog.stats/pb/elog"
	"github.com/spf13/cobra"
)

var (
	cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create new event",
		Long:  "Create new event",
		Run: func(cmd *cobra.Command, args []string) {
			execute(createEvent)
		},
	}

	cmdList = &cobra.Command{
		Use:   "list",
		Short: "List events with filters",
		Long:  `List events with filters. If filters is omit, return all events.`,
		Run: func(cmd *cobra.Command, args []string) {
			execute(listEvents)
		},
	}

	rootCmd = &cobra.Command{Use: "elog"}

	host     string
	port     string
	ipserver string
	ipclient string
	message  string
	tags     []string
)

func init() {
	cmdCreate.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "server address/ip")
	cmdCreate.Flags().StringVarP(&port, "port", "p", "50051", "server port")
	cmdCreate.Flags().StringVarP(&message, "msg", "m", "", "event.Message")
	cmdCreate.Flags().StringVarP(&ipserver, "ipserver", "s", "", "event.IPServer")
	cmdCreate.Flags().StringVarP(&ipclient, "ipclient", "c", "", "event.IPClient")
	cmdCreate.Flags().StringArrayVarP(&tags, "tags", "t", []string{}, "event.Tags \"key:value\"")

	cmdList.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "server address/ip")
	cmdList.Flags().StringVarP(&port, "port", "p", "50051", "server port")
	cmdList.Flags().StringVarP(&ipserver, "ipserver", "s", "", "filter.IPServer")
	cmdList.Flags().StringVarP(&ipclient, "ipclient", "c", "", "filter.IPClient")
	cmdList.Flags().StringArrayVarP(&tags, "tags", "t", []string{}, "filter.Tags \"key:value\"")

	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdList)
}

func main() {
	rootCmd.Execute()
}

func execute(f func(context.Context, pb.ElogClient)) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc: unable to connect to server, %v\n", err)
	}
	defer conn.Close()

	client := pb.NewElogClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	f(ctx, client)
}

func createEvent(ctx context.Context, client pb.ElogClient) {
	mTags, err := convertToMapTags(tags)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	pbEvents := &pb.Event{
		IpClient: ipclient,
		IpServer: ipserver,
		Message:  message,
		Tags:     mTags,
	}
	_, err = client.Create(ctx, pbEvents)
	if err != nil {
		log.Fatalf("elog: unable to create event, %v\n", err)
	}
	fmt.Printf(`==> DONE
   event:
   - ipclient: %s
   - ipserver: %s
   - message:  %s
   - tags:     %v
`, ipclient, ipserver, message, mTags)
}

func listEvents(ctx context.Context, client pb.ElogClient) {
	mTags, err := convertToMapTags(tags)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	pbQuery := &pb.Query{
		IpClient: ipclient,
		IpServer: ipserver,
		Tags:     mTags,
	}
	stream, err := client.List(ctx, pbQuery)
	if err != nil {
		log.Fatalf("elog: unable to list events, %v\n", err)
	}
	count := 0
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("elog: unable to list events, %v\n", err)
		}
		count = count + 1
		fmt.Printf("%d. %v\n", count, event)
	}
	fmt.Printf("Total: %d events\n", count)
}

func convertToMapTags(values []string) (map[string]string, error) {
	var mTags = map[string]string{}
	for _, value := range values {
		items := strings.Split(value, ":")
		if len(items) != 2 {
			return nil, fmt.Errorf("tags: unknown format, required key:value")
		}
		mTags[items[0]] = items[1]
	}
	return mTags, nil
}
