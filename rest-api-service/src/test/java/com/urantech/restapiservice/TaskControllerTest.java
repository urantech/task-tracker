package com.urantech.restapiservice;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.urantech.restapiservice.model.entity.User;
import com.urantech.restapiservice.model.entity.UserAuthority;
import com.urantech.restapiservice.model.rest.TaskDto;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.test.web.servlet.request.SecurityMockMvcRequestPostProcessors;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.transaction.annotation.Propagation;
import org.springframework.transaction.annotation.Transactional;

import java.util.Set;

import static com.urantech.restapiservice.model.entity.UserAuthority.Authority.USER;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@ActiveProfiles("test")
@AutoConfigureMockMvc(printOnlyOnFailure = false)
class TaskControllerTest {
    @Autowired
    private MockMvc mockMvc;
    @Autowired
    private ObjectMapper mapper;
    @Autowired
    private JdbcTemplate jdbcTemplate;

    private final User user = new User();
    private final UserAuthority authority = new UserAuthority();
    private final TaskDto taskDto = new TaskDto(1L, "task", false);

    @BeforeEach
    void setUp() {
        authority.setAuthority(USER);
        setUpSecurityContext();

        jdbcTemplate.execute("truncate table users, user_authority, task restart identity cascade");
        jdbcTemplate.execute("insert into users (id, email, password, enabled) "
                        + "values (12345, 'j.dewar@gmail.com', 'some pass', true)");
        jdbcTemplate.execute("insert into task (id, description, user_id, done) "
                        + "values (12345, 'task1', 12345, false), (12346, 'task2', 12345, true)");
    }

    @Test
    void shouldGetCurrentUserTasks() throws Exception {
        // given
        var requestBuilder =
                MockMvcRequestBuilders.get("/api/tasks")
                        .with(
                                SecurityMockMvcRequestPostProcessors.securityContext(
                                        SecurityContextHolder.getContext()));

        // when
        mockMvc.perform(requestBuilder)
                // then
                .andExpectAll(
                        status().isOk(),
                        jsonPath("$.length()").value(2),
                        jsonPath("$[0].description").value("task1"),
                        jsonPath("$[0].done").value(false),
                        jsonPath("$[1].description").value("task2"),
                        jsonPath("$[1].done").value(true));
    }

    @Test
    void shouldCreateNewTask() throws Exception {
        // given
        var requestBuilder =
                MockMvcRequestBuilders.post("/api/tasks")
                        .with(
                                SecurityMockMvcRequestPostProcessors.securityContext(
                                        SecurityContextHolder.getContext()))
                        .content(mapper.writeValueAsString(taskDto))
                        .contentType(MediaType.APPLICATION_JSON);

        // when
        mockMvc.perform(requestBuilder)
                // then
                .andExpectAll(
                        status().isOk(),
                        jsonPath("$.length()").value(3),
                        jsonPath("$.description").value("task"),
                        jsonPath("$.done").value(false));
    }

    @Test
    void shouldUpdateTask() throws Exception {
        // given
        TaskDto toUpdate = new TaskDto(12345L, "updated", true);
        var requestBuilder =
                MockMvcRequestBuilders.patch("/api/tasks")
                        .with(
                                SecurityMockMvcRequestPostProcessors.securityContext(
                                        SecurityContextHolder.getContext()))
                        .content(mapper.writeValueAsString(toUpdate))
                        .contentType(MediaType.APPLICATION_JSON);

        // when
        mockMvc.perform(requestBuilder)
                // then
                .andExpectAll(
                        status().isOk(),
                        jsonPath("$.length()").value(3),
                        jsonPath("$.description").value("updated"),
                        jsonPath("$.done").value(true));
    }

    @Test
    @Transactional(propagation = Propagation.NOT_SUPPORTED)
    void shouldDeleteTask() throws Exception {
        // given
        var requestBuilder =
                MockMvcRequestBuilders.delete("/api/tasks/{id}", 12345)
                        .with(
                                SecurityMockMvcRequestPostProcessors.securityContext(
                                        SecurityContextHolder.getContext()));

        // when
        mockMvc.perform(requestBuilder).andExpect(status().isOk());

        // then
        Integer actual =
                jdbcTemplate.queryForObject(
                        "select count(*) from task where id = 12345", Integer.class);
        Assertions.assertEquals(0, actual);
    }

    private void setUpSecurityContext() {
        user.setId(12345L);
        user.setEmail("j.dewar@gmail.com");
        Authentication auth =
                new UsernamePasswordAuthenticationToken(user, null, Set.of(authority));
        SecurityContextHolder.getContext().setAuthentication(auth);
    }
}
