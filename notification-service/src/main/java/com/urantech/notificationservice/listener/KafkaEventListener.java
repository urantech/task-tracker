package com.urantech.notificationservice.listener;

import com.urantech.notificationservice.model.EmailMessage;
import com.urantech.notificationservice.service.EmailService;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@RequiredArgsConstructor
public class KafkaEventListener {

    private final EmailService emailService;

    @KafkaListener(
            topics = "user_registration",
            groupId = "rest-api-service"
    )
    public void handleUserRegistrationEvent(String message) {
        emailService.sendEmail(EmailMessage.build(message));
    }

    @KafkaListener(
            topics = "daily_report",
            groupId = "report-service"
    )
    public void handleDailyReportEvent(String message) {
        emailService.sendEmail(EmailMessage.build(message));
    }
}
