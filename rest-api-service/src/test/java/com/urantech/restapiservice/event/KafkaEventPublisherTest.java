package com.urantech.restapiservice.event;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.bean.override.mockito.MockitoBean;

import static org.mockito.Mockito.verify;

@SpringBootTest
@ActiveProfiles("test")
class KafkaEventPublisherTest {
    @Autowired
    private KafkaEventPublisher publisher;
    @MockitoBean
    private KafkaTemplate<String, String> kafkaTemplate;

    @Test
    void shouldSendRegistrationEvent() {
        String testEmail = "test@example.com";

        publisher.publishUserRegistrationEvent(testEmail);

        verify(kafkaTemplate).send("user_registration", testEmail);
    }
}
