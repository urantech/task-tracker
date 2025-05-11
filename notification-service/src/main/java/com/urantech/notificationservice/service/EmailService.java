package com.urantech.notificationservice.service;

import com.urantech.notificationservice.model.EmailMessage;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class EmailService {

    private final JavaMailSender emailSender;

    public void sendEmail(EmailMessage emailMessage) {
        SimpleMailMessage msg = new SimpleMailMessage();
        msg.setTo(emailMessage.getTo());
        msg.setSubject(emailMessage.getSubject());
        msg.setText(emailMessage.getBody());

        emailSender.send(msg);
        log.info("Email sent to {} successfully", emailMessage.getTo());
    }
}
