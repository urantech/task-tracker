package com.urantech.restapiservice.service.auth;

import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.oauth2.jwt.JwsHeader;
import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.jwt.JwtEncoderParameters;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Duration;
import java.time.Instant;

@Service
@RequiredArgsConstructor
public class JwtService {

    private final JwtEncoder jwtEncoder;

    @Transactional
    public String generateToken(Authentication auth) {
        String email = auth.getName();
        Instant now = Instant.now();

        JwtClaimsSet claims = JwtClaimsSet.builder()
                .claim("authorities",
                        auth.getAuthorities().stream().map(GrantedAuthority::getAuthority))
                .issuedAt(now)
                .expiresAt(now.plus(Duration.ofDays(1)))
                .subject(email)
                .build();

        JwsHeader jwsHeader = JwsHeader.with(() -> "HS256").build();
        return jwtEncoder.encode(JwtEncoderParameters.from(jwsHeader, claims)).getTokenValue();
    }
}
