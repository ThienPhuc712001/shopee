package email

import (
	"bytes"
	"crypto/tls"
	"ecommerce/internal/domain/model"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

// EmailService handles email sending via SMTP
type EmailService struct {
	config model.EmailConfig
}

// NewEmailService creates a new email service
func NewEmailService(config model.EmailConfig) *EmailService {
	return &EmailService{config: config}
}

// SendEmail sends an email using SMTP
func (s *EmailService) SendEmail(req model.EmailRequest) error {
	// Build email message
	message := s.buildMessage(req)

	// Set up authentication
	var auth smtp.Auth
	if s.config.Username != "" {
		auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	// Build server address
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Send email
	if s.config.UseTLS {
		return s.sendWithTLS(addr, auth, req.To, message)
	}
	
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	return smtp.SendMail(addr, auth, from, req.To, []byte(message))
}

// sendWithTLS sends email using TLS connection
func (s *EmailService) sendWithTLS(addr string, auth smtp.Auth, to []string, message string) error {
	// Set up TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}

	// Connect to server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	// Authenticate if credentials provided
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	// Set sender
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	if err := client.Mail(from); err != nil {
		return err
	}

	// Set recipients
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	// Get writer
	writer, err := client.Data()
	if err != nil {
		return err
	}
	defer writer.Close()

	// Write message
	_, err = writer.Write([]byte(message))
	return err
}

// buildMessage builds the email message string
func (s *EmailService) buildMessage(req model.EmailRequest) string {
	var buf bytes.Buffer

	// From header
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))

	// To header
	buf.WriteString(fmt.Sprintf("To: %s\r\n", req.ToName))

	// Subject header
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))

	// Date header
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))

	// MIME header
	if req.HTML {
		buf.WriteString("MIME-Version: 1.0\r\n")
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	buf.WriteString("\r\n")

	// Body
	buf.WriteString(req.Body)

	return buf.String()
}

// SendOrderConfirmation sends order confirmation email
func (s *EmailService) SendOrderConfirmation(toEmail, toName, orderNumber string, total float64) error {
	subject := fmt.Sprintf("Order Confirmation - %s", orderNumber)
	body := s.renderTemplate(orderConfirmationTemplate, map[string]string{
		"OrderNumber": orderNumber,
		"Total":       fmt.Sprintf("$%.2f", total),
		"CustomerName": toName,
	})

	return s.SendEmail(model.EmailRequest{
		To:       []string{toEmail},
		ToName:   toName,
		Subject:  subject,
		Body:     body,
		HTML:     true,
	})
}

// SendPaymentConfirmation sends payment receipt email
func (s *EmailService) SendPaymentConfirmation(toEmail, toName, orderNumber string, amount float64, paymentMethod string) error {
	subject := fmt.Sprintf("Payment Received - %s", orderNumber)
	body := s.renderTemplate(paymentConfirmationTemplate, map[string]string{
		"OrderNumber":   orderNumber,
		"Amount":        fmt.Sprintf("$%.2f", amount),
		"CustomerName":  toName,
		"PaymentMethod": paymentMethod,
	})

	return s.SendEmail(model.EmailRequest{
		To:       []string{toEmail},
		ToName:   toName,
		Subject:  subject,
		Body:     body,
		HTML:     true,
	})
}

// SendShippingUpdate sends shipping notification email
func (s *EmailService) SendShippingUpdate(toEmail, toName, orderNumber, trackingNumber, status string) error {
	subject := fmt.Sprintf("Shipping Update - %s", orderNumber)
	body := s.renderTemplate(shippingUpdateTemplate, map[string]string{
		"OrderNumber":    orderNumber,
		"TrackingNumber": trackingNumber,
		"Status":         status,
		"CustomerName":   toName,
	})

	return s.SendEmail(model.EmailRequest{
		To:       []string{toEmail},
		ToName:   toName,
		Subject:  subject,
		Body:     body,
		HTML:     true,
	})
}

// SendPromotionEmail sends promotional email
func (s *EmailService) SendPromotionEmail(toEmail, toName, promoTitle, promoCode string, discountPercent float64) error {
	subject := fmt.Sprintf("Special Offer: %s", promoTitle)
	body := s.renderTemplate(promotionTemplate, map[string]string{
		"PromoTitle":      promoTitle,
		"PromoCode":       promoCode,
		"DiscountPercent": fmt.Sprintf("%.0f", discountPercent),
		"CustomerName":    toName,
	})

	return s.SendEmail(model.EmailRequest{
		To:       []string{toEmail},
		ToName:   toName,
		Subject:  subject,
		Body:     body,
		HTML:     true,
	})
}

// SendPasswordReset sends password reset email
func (s *EmailService) SendPasswordReset(toEmail, toName, resetLink string) error {
	subject := "Password Reset Request"
	body := s.renderTemplate(passwordResetTemplate, map[string]string{
		"ResetLink":    resetLink,
		"CustomerName": toName,
	})

	return s.SendEmail(model.EmailRequest{
		To:       []string{toEmail},
		ToName:   toName,
		Subject:  subject,
		Body:     body,
		HTML:     true,
	})
}

// SendWelcomeEmail sends welcome email to new users
func (s *EmailService) SendWelcomeEmail(toEmail, toName string) error {
	subject := "Welcome to Our Store!"
	body := s.renderTemplate(welcomeTemplate, map[string]string{
		"CustomerName": toName,
	})

	return s.SendEmail(model.EmailRequest{
		To:       []string{toEmail},
		ToName:   toName,
		Subject:  subject,
		Body:     body,
		HTML:     true,
	})
}

// renderTemplate renders an HTML template with data
func (s *EmailService) renderTemplate(templateStr string, data map[string]string) string {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return templateStr
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return templateStr
	}

	return buf.String()
}

// Email Templates
const orderConfirmationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .order-details { background: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; padding: 10px 20px; background: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Order Confirmation</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>Thank you for your order! We've received your order and are processing it.</p>
            
            <div class="order-details">
                <h3>Order Details</h3>
                <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
                <p><strong>Total Amount:</strong> {{.Total}}</p>
            </div>
            
            <p>We'll send you another email when your order ships.</p>
            
            <p style="text-align: center; margin: 20px 0;">
                <a href="#" class="button">Track Your Order</a>
            </p>
        </div>
        <div class="footer">
            <p>© 2024 E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

const paymentConfirmationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #2196F3; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .payment-details { background: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Payment Received</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>Your payment has been successfully processed.</p>
            
            <div class="payment-details">
                <h3>Payment Details</h3>
                <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
                <p><strong>Amount Paid:</strong> {{.Amount}}</p>
                <p><strong>Payment Method:</strong> {{.PaymentMethod}}</p>
            </div>
            
            <p>Thank you for your purchase!</p>
        </div>
        <div class="footer">
            <p>© 2024 E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

const shippingUpdateTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #FF9800; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .shipping-details { background: white; padding: 15px; margin: 15px 0; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; padding: 10px 20px; background: #FF9800; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Shipping Update</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>Good news! Your order has been shipped.</p>
            
            <div class="shipping-details">
                <h3>Shipping Information</h3>
                <p><strong>Order Number:</strong> {{.OrderNumber}}</p>
                <p><strong>Tracking Number:</strong> {{.TrackingNumber}}</p>
                <p><strong>Status:</strong> {{.Status}}</p>
            </div>
            
            <p style="text-align: center; margin: 20px 0;">
                <a href="#" class="button">Track Your Package</a>
            </p>
        </div>
        <div class="footer">
            <p>© 2024 E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

const promotionTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #E91E63; color: white; padding: 30px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .promo-box { background: white; padding: 20px; margin: 15px 0; border-radius: 5px; text-align: center; border: 2px dashed #E91E63; }
        .promo-code { font-size: 24px; font-weight: bold; color: #E91E63; letter-spacing: 2px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; padding: 15px 30px; background: #E91E63; color: white; text-decoration: none; border-radius: 5px; font-size: 16px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Special Offer!</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>We have a special offer just for you!</p>
            
            <div class="promo-box">
                <h2>{{.PromoTitle}}</h2>
                <p>Get <strong>{{.DiscountPercent}}% OFF</strong> your next order</p>
                <p>Use code:</p>
                <p class="promo-code">{{.PromoCode}}</p>
            </div>
            
            <p style="text-align: center; margin: 20px 0;">
                <a href="#" class="button">Shop Now</a>
            </p>
            
            <p><small>This offer is valid for a limited time only. Don't miss out!</small></p>
        </div>
        <div class="footer">
            <p>© 2024 E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

const passwordResetTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #607D8B; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .button { display: inline-block; padding: 15px 30px; background: #607D8B; color: white; text-decoration: none; border-radius: 5px; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .warning { background: #fff3cd; padding: 15px; border-left: 4px solid #ffc107; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>We received a request to reset your password. Click the button below to reset it:</p>
            
            <p style="text-align: center; margin: 20px 0;">
                <a href="{{.ResetLink}}" class="button">Reset Password</a>
            </p>
            
            <div class="warning">
                <strong>Important:</strong> This link will expire in 1 hour. If you didn't request this, please ignore this email.
            </div>
        </div>
        <div class="footer">
            <p>© 2024 E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .welcome-box { background: white; padding: 20px; margin: 15px 0; border-radius: 5px; text-align: center; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; padding: 15px 30px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome!</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            <p>Welcome to our store! We're excited to have you as part of our community.</p>
            
            <div class="welcome-box">
                <h3>Get Started</h3>
                <p>Browse our collection and discover amazing products.</p>
                <p>As a welcome gift, use code <strong>WELCOME10</strong> for 10% off your first order!</p>
            </div>
            
            <p style="text-align: center; margin: 20px 0;">
                <a href="#" class="button">Start Shopping</a>
            </p>
            
            <p>If you have any questions, feel free to contact our support team.</p>
        </div>
        <div class="footer">
            <p>© 2024 E-Commerce Store. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
