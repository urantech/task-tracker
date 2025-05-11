package com.urantech.reportservice.client;

import com.urantech.reportservice.model.rest.UserResponse;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.List;

@FeignClient(name = "rest-api-service")
public interface UserClient {

    @GetMapping("/api/admin/users")
    List<UserResponse> getAllUsersWithUnfinishedTasks();

    // todo: add jwt auth
}
