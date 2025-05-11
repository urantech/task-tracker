package com.urantech.restapiservice.service.user;

import com.urantech.restapiservice.event.KafkaEventPublisher;
import com.urantech.restapiservice.model.entity.User;
import com.urantech.restapiservice.model.entity.UserAuthority;
import com.urantech.restapiservice.model.rest.user.RegistrationRequest;
import com.urantech.restapiservice.repository.UserAuthorityRepository;
import com.urantech.restapiservice.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.util.Set;

@Service
@RequiredArgsConstructor
public class UserService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final KafkaEventPublisher eventPublisher;
    private final UserAuthorityRepository authorityRepository;

    @Transactional
    public void register(RegistrationRequest req) {
        User user = new User(req.email(), passwordEncoder.encode(req.password()));
        UserAuthority authority = new UserAuthority(UserAuthority.Authority.USER, user);
        user.setAuthorities(Set.of(authority));
        try {
            userRepository.save(user);
            authorityRepository.save(authority);
        } catch (DataIntegrityViolationException e) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "User already exists");
        }
        eventPublisher.publishUserRegistrationEvent(req.email());
    }
}
