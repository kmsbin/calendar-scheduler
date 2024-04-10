package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"html/template"
	"log"
	"os"
)

const (
	Sender  = "kaulisabin@gmail.com"
	Subject = "Password recovery"
	CharSet = "UTF-8"
)

type SESService struct{}

func NewSESService() SESService {
	return SESService{}
}

func (e SESService) SendEmail(emailData EmailData) error {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := ses.New(sess)
	fullUrl := fmt.Sprintf("%s/app/auth/reset-password/%s", emailData.BaseUrl, emailData.Code)
	emailTemplateParsed, err := getTemplateEmail(fullUrl)
	if err != nil {
		return err
	}

	input := getSESEmailInput(emailData.Email, emailTemplateParsed)
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}

		return err
	}

	log.Println("Email Sent to address: " + emailData.Email)
	log.Printf("Result %v\n", result.String())
	return nil
}

func getSESEmailInput(destinationEmail, emailTemplate string) *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(destinationEmail)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(emailTemplate),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(emailTemplate),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}
}

func getTemplateEmail(link string) (string, error) {
	emailTemplateFile, err := os.ReadFile("public/email_template.html")

	if err != nil {
		log.Println(err)
		return "", errors.New("err when open email template")
	}
	emailTemplate, err := template.New("").Parse(string(emailTemplateFile))
	if err != nil {
		return "", errors.New("error when try to parse email template")
	}
	var tpl bytes.Buffer
	data := map[string]string{"link": link}
	err = emailTemplate.Execute(&tpl, data)

	if err != nil {
		return "", errors.New("error executing template")
	}
	return tpl.String(), nil
}
