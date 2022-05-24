package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/* template methods that can be varied even though the core algo remains the same */
type iAuthTemplate interface {
	fetchDestination(userId string) string
	sendChallenge(dest string)
	persist(userId string, otp int)
	validate(userId string, userInput int) bool
}

/* handle to invoke the tenplate methods of specific impl */
type SimpleOtpBasedAuthentication struct {
	auth iAuthTemplate
}

/* core algo */
func (otp *SimpleOtpBasedAuthentication) authenticate(userId string) {
	dest := otp.auth.fetchDestination(userId)
	otp.auth.sendChallenge(dest)
	otpNum := rand.Intn(9999)
	otp.auth.persist(userId, otpNum)
	//spin go routine to receive otp from user on chan
	otp.auth.validate(userId, otpNum)
}

/* one type opf impl for the template */
type SMS struct {
	once  sync.Once
	cache map[string]int
}

func (smsAuth *SMS) fetchDestination(id string) string {
	return "9880529612"
}

func (smsAuth *SMS) validate(id string, otp int) bool {
	fmt.Println("Validating otp of user ", otp, " ", id)
	if _, ok := smsAuth.cache[id]; ok {
		delete(smsAuth.cache, id)
		return true
	}
	return false

}

func (smsAuth *SMS) sendChallenge(phoneNo string) {
	fmt.Println("Send otp to phone number", phoneNo, " ", rand.Intn(9999))
}

func (smsAuth *SMS) persist(userId string, otp int) {
	fmt.Println("Caching otp")
	smsAuth.once.Do(func() { smsAuth.cache = make(map[string]int) })
	smsAuth.cache[userId] = otp
}

/* another type opf impl for the template */

type EMail struct {
}

func (emailAuth *EMail) sendChallenge(emailAddr string) {
	fmt.Println("Send otp to email ", rand.Intn(9999))
}

func (emailAuth *EMail) persist(userId string, otp int) {
	fmt.Println("Caching otp")
}

func (emailAuth *EMail) fetchDestination(userId string) string {
	return ""
}

func (emailAuth *EMail) validate(id string, otp int) bool {
	fmt.Println("Validating otp of user ", otp, " ", id)
	return false
}

func main() {
	sms := &SMS{}
	simpleAuth1 := &SimpleOtpBasedAuthentication{auth: sms}
	simpleAuth1.authenticate("shashank")
	fmt.Println("------------------------------------------------------------")
	email := &EMail{}
	simpleAuth2 := &SimpleOtpBasedAuthentication{auth: email}
	simpleAuth2.authenticate("shashank")

}
