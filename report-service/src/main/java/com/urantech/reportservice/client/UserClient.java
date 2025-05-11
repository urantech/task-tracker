package com.urantech.reportservice.client;

import com.urantech.reportservice.config.StaticJwtFeignConfig;
import com.urantech.reportservice.model.rest.UserDto;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.List;

@FeignClient(
        name = "rest-api-service",
        configuration = StaticJwtFeignConfig.class
)
public interface UserClient {

    @GetMapping("/api/users/unfinishedTasks")
    List<UserDto> getAllUsersWithUnfinishedTasks();
}
