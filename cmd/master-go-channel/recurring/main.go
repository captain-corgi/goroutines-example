package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/captain-corgi/goroutines-example/pkg/recurring"
)

type (
	PendingUserNotifications map[int][]*Notification
	Notification             struct {
		Content string
		UserId  int
	}
)

func main() {
	pendingNotificationsProcess()
}

func pendingNotificationsProcess() {
	process := &recurring.RecurringProcess{}
	notifications := PendingUserNotifications{}
	handler := func() {
		collectNewUsersNotifications(notifications)
		handlePendingUsersNotifications(notifications, sendUserBatchNotificationsEmail, process)
	}
	interval := 1 * time.Second
	startTime := time.Now().Add(5 * time.Second)
	process = createRecurringProcess("Pending User Notifications", handler, interval, startTime)

	<-process.Stop()
}

func sendUserBatchNotificationsEmail(userId int, notifications []*Notification) {
	fmt.Printf("Sending email to user with userId %d for pending notifications %v\n", userId, notifications)
}

func handlePendingUsersNotifications(pendingNotifications PendingUserNotifications, handler func(userId int, notifications []*Notification), process *recurring.RecurringProcess) {
	userNotificationCount := 0
	for userId, notifications := range pendingNotifications {
		userNotificationCount++
		handler(userId, notifications)
		delete(pendingNotifications, userId)
	}

	if userNotificationCount == 0 {
		process.Cancel()
	}
}

func collectNewUsersNotifications(notifications PendingUserNotifications) {
	randomNotifications := getRandomNotifications()
	if len(randomNotifications) > 0 {
		notifications[randomNotifications[0].UserId] = randomNotifications
	}
}

func getRandomNotifications() (notifications []*Notification) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	userId := rand.Intn(100-10+1) + 10
	numOfNotifications := rand.Intn(5-0+1) + 0
	fmt.Printf("numOfNotifications %+v\n", numOfNotifications)
	for i := 0; i < numOfNotifications; i++ {
		notifications = append(notifications, &Notification{Content: gofakeit.Paragraph(1, 2, 10, " "), UserId: userId})
	}

	return
}

func createRecurringProcess(name string, handler recurring.ProcessHandler, interval time.Duration, startTime time.Time) *recurring.RecurringProcess {
	process := recurring.NewRecurringProcess(name, interval, startTime, handler, make(chan struct{}))
	go process.Start()
	return process
}
