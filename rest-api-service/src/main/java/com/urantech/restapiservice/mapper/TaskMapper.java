package com.urantech.restapiservice.mapper;

import com.urantech.restapiservice.model.entity.Task;
import com.urantech.restapiservice.model.rest.TaskDto;
import org.mapstruct.Mapper;
import org.mapstruct.MappingConstants;
import org.mapstruct.MappingTarget;
import org.mapstruct.ReportingPolicy;

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE, componentModel = MappingConstants.ComponentModel.SPRING)
public interface TaskMapper {
    Task toEntity(TaskDto taskDto);

    TaskDto toTaskDto(Task task);

    Task updateWithNull(TaskDto taskDto, @MappingTarget Task task);
}
