package com.urantech.restapiservice.service;

import com.urantech.restapiservice.model.entity.Task;
import com.urantech.restapiservice.model.entity.User;
import com.urantech.restapiservice.model.rest.TaskDto;
import com.urantech.restapiservice.repository.TaskRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class TaskService {

    private final TaskRepository taskRepository;

    public TaskDto create(TaskDto taskDto, User user) {
        Task task = new Task(taskDto.description(), user);
        Task entity = taskRepository.save(task);
        return TaskDto.fromEntity(entity);
    }

    public List<TaskDto> getTasks(User user) {
        return taskRepository.findAllByUserId(user.getId())
                .stream()
                .map(TaskDto::fromEntity)
                .toList();
    }

    public TaskDto update(TaskDto taskDto) {
        Optional<Task> taskOpt = taskRepository.findById(taskDto.id());
        if (taskOpt.isEmpty()) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Task with id: %d not found".formatted(taskDto.id()));
        }

        Task task = taskOpt.get();
        task.setDescription(taskDto.description());
        task.setDone(taskDto.done());

        Task entity = taskRepository.save(task);
        return TaskDto.fromEntity(entity);
    }

    public void deleteById(long id) {
        taskRepository.deleteById(id);
    }
}
