package com.urantech.restapiservice.model.rest;

import com.urantech.restapiservice.model.entity.Task;

public record TaskDto(Long id, String description, boolean done) {

    public static TaskDto fromEntity(Task task) {
        return new TaskDto(task.getId(), task.getDescription(), task.isDone());
    }
}
