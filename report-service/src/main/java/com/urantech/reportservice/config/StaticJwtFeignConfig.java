package com.urantech.reportservice.config;

import feign.RequestInterceptor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class StaticJwtFeignConfig {

    @Value("${rest-api-service.auth.jwt-token}")
    private String jwtToken;

    @Bean
    public RequestInterceptor jwtInterceptor() {
        return requestTemplate -> requestTemplate.header("Authorization", "Bearer " + jwtToken);
    }
}
