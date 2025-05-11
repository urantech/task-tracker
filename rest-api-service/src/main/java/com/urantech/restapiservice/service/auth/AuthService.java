package com.urantech.restapiservice.service.auth;

import com.urantech.restapiservice.model.rest.auth.AuthRequest;
import com.urantech.restapiservice.model.rest.auth.AuthResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class AuthService {

    private final JwtService jwtService;
    private final AuthenticationManager authManager;

    public AuthResponse authenticate(AuthRequest req) {
        String email = req.email().trim().toLowerCase();
        Authentication auth = authManager.authenticate(new UsernamePasswordAuthenticationToken(email, req.password()));
        String token = jwtService.generateToken(auth);
        return new AuthResponse(token);
    }
}
