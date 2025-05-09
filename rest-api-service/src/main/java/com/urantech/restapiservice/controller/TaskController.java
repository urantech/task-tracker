package com.urantech.restapiservice.controller;
import com.fasterxml.jackson.databind.JsonNode;
import com.urantech.restapiservice.model.rest.TaskDto;
import com.urantech.restapiservice.service.TaskService;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PagedModel;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.List;

@RestController
@RequiredArgsConstructor
@RequestMapping("/api/tasks")
public class TaskController {

    private final TaskService taskService;

    @GetMapping
    public PagedModel<TaskDto> getAll(Pageable pageable) {
        Page<TaskDto> taskDtos = taskService.getAll(pageable);
        return new PagedModel<>(taskDtos);
    }

    @GetMapping("/{id}")
    public TaskDto getOne(@PathVariable Long id) {
        return taskService.getOne(id);
    }

    @GetMapping("/by-ids")
    public List<TaskDto> getMany(@RequestParam List<Long> ids) {
        return taskService.getMany(ids);
    }

    @PostMapping
    public TaskDto create(@RequestBody TaskDto dto) {
        return taskService.create(dto);
    }

    @PatchMapping("/{id}")
    public TaskDto patch(@PathVariable Long id, @RequestBody JsonNode patchNode) throws IOException {
        return taskService.patch(id, patchNode);
    }

    @PatchMapping
    public List<Long> patchMany(@RequestParam List<Long> ids, @RequestBody JsonNode patchNode) throws IOException {
        return taskService.patchMany(ids, patchNode);
    }

    @DeleteMapping("/{id}")
    public TaskDto delete(@PathVariable Long id) {
        return taskService.delete(id);
    }

    @DeleteMapping
    public void deleteMany(@RequestParam List<Long> ids) {
        taskService.deleteMany(ids);
    }
}
