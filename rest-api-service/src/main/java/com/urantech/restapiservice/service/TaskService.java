package com.urantech.restapiservice.service;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.urantech.restapiservice.mapper.TaskMapper;
import com.urantech.restapiservice.model.entity.Task;
import com.urantech.restapiservice.model.rest.TaskDto;
import com.urantech.restapiservice.repository.TaskRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.io.IOException;
import java.util.Collection;
import java.util.List;
import java.util.Optional;

@Service
@RequiredArgsConstructor
public class TaskService{

    private final TaskMapper taskMapper;
    private final ObjectMapper objectMapper;
    private final TaskRepository taskRepository;

    public Page<TaskDto> getAll(Pageable pageable) {
        Page<Task> tasks = taskRepository.findAll(pageable);
        return tasks.map(taskMapper::toTaskDto);
    }

    public TaskDto getOne(Long id) {
        Optional<Task> taskOptional = taskRepository.findById(id);
        return taskMapper.toTaskDto(taskOptional.orElseThrow(() ->
                new ResponseStatusException(HttpStatus.NOT_FOUND, "Entity with id `%s` not found".formatted(id))));
    }

    public List<TaskDto> getMany(List<Long> ids) {
        List<Task> tasks = taskRepository.findAllById(ids);
        return tasks.stream()
                .map(taskMapper::toTaskDto)
                .toList();
    }

    public TaskDto create(TaskDto dto) {
        if (dto.id() != null) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "id must be null");
        }
        Task task = taskMapper.toEntity(dto);
        Task resultTask = taskRepository.save(task);
        return taskMapper.toTaskDto(resultTask);
    }

    public TaskDto patch(Long id, JsonNode patchNode) throws IOException {
        Task task = taskRepository.findById(id).orElseThrow(() ->
                new ResponseStatusException(HttpStatus.NOT_FOUND, "Entity with id `%s` not found".formatted(id)));

        TaskDto taskDto = taskMapper.toTaskDto(task);
        objectMapper.readerForUpdating(taskDto).readValue(patchNode);
        taskMapper.updateWithNull(taskDto, task);

        Task resultTask = taskRepository.save(task);
        return taskMapper.toTaskDto(resultTask);
    }

    public List<Long> patchMany(List<Long> ids, JsonNode patchNode) throws IOException {
        Collection<Task> tasks = taskRepository.findAllById(ids);

        for (Task task : tasks) {
            TaskDto taskDto = taskMapper.toTaskDto(task);
            objectMapper.readerForUpdating(taskDto).readValue(patchNode);
            taskMapper.updateWithNull(taskDto, task);
        }

        List<Task> resultTasks = taskRepository.saveAll(tasks);
        return resultTasks.stream()
                .map(Task::getId)
                .toList();
    }

    public TaskDto delete(Long id) {
        Task task = taskRepository.findById(id).orElse(null);
        if (task != null) {
            taskRepository.delete(task);
        }
        return taskMapper.toTaskDto(task);
    }

    public void deleteMany(List<Long> ids) {
        taskRepository.deleteAllById(ids);
    }
}
