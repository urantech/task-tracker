package com.urantech.restapiservice.service.user;

import com.urantech.restapiservice.model.entity.UserAuthority;
import com.urantech.restapiservice.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Set;

@Service
@RequiredArgsConstructor
public class CustomUserDetailsService implements UserDetailsService {

    private final UserRepository userRepository;

    @Override
    public UserDetails loadUserByUsername(String email) throws UsernameNotFoundException {
        return userRepository.findByEmailWithAuthorities(email)
                .map(CustomUserDetailsService::mapToUserDetails)
                .orElseThrow(() -> new UsernameNotFoundException("User with email %s not found".formatted(email)));
    }

    private static UserDetails mapToUserDetails(com.urantech.restapiservice.model.entity.User user) {
        return User.builder()
                .username(user.getEmail())
                .password(user.getPassword())
                .authorities(mapToSimpleGrantedAuthorities(user.getAuthorities()))
                .build();
    }

    private static List<SimpleGrantedAuthority> mapToSimpleGrantedAuthorities(Set<UserAuthority> authorities) {
        return authorities.stream()
                .map(UserAuthority::getAuthority)
                .map(SimpleGrantedAuthority::new)
                .toList();
    }
}
