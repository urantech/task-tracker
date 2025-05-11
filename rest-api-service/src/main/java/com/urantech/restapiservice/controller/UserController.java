package com.urantech.restapiservice.controller;

import com.urantech.restapiservice.model.entity.User;
import com.urantech.restapiservice.model.rest.UserDto;
import com.urantech.restapiservice.model.rest.user.RegistrationRequest;
import com.urantech.restapiservice.model.rest.user.UserResponse;
import com.urantech.restapiservice.service.user.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/users")
public class UserController {

    private final UserService userService;

    @PostMapping("/register")
    public void register(@RequestBody RegistrationRequest req) {
        userService.register(req);
    }

    @GetMapping("/unfinishedTasks")
    public List<UserDto> getUsersWithUnfinishedTasks() {
        return userService.getUsersWithUnfinishedTasks();
    }

    @GetMapping("/user")
    public UserResponse getUser(@AuthenticationPrincipal User user) {
        return new UserResponse(user.getEmail());
    }
}
