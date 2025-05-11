package com.urantech.restapiservice.repository;

import com.urantech.restapiservice.model.entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {

    @Query("select u from User u join fetch u.authorities where u.email = :email")
    Optional<User> findByEmailWithAuthorities(String email);
}
