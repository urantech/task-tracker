package com.urantech.notificationservice.model;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class EmailMessage {
    private String to;
    private String subject;
    private String body;

    public static EmailMessage build(String message) {
        return EmailMessage.builder()
                .to(message)
                .subject(message)
                .body("Registration successful!")
                .build();
    }
}
