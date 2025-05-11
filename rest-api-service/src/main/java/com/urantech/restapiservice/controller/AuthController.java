package com.urantech.restapiservice.controller;

import com.urantech.restapiservice.model.rest.auth.AuthRequest;
import com.urantech.restapiservice.model.rest.auth.AuthResponse;
import com.urantech.restapiservice.service.auth.AuthService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/auth")
public class AuthController {

    private final AuthService authService;

    @PostMapping("/login")
    public AuthResponse login(@RequestBody AuthRequest req) {
        return authService.authenticate(req);
    }
}
