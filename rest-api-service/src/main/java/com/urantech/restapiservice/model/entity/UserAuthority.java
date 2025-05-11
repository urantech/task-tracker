package com.urantech.restapiservice.model.entity;

import jakarta.persistence.*;
import lombok.AccessLevel;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.proxy.HibernateProxy;
import org.springframework.security.core.GrantedAuthority;

import java.util.Objects;

@Getter
@Setter
@Entity
@NoArgsConstructor
@Table(name = "user_authority")
public class UserAuthority implements GrantedAuthority {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false)
    private Long id;

    @Column(name = "authority")
    @Enumerated(EnumType.STRING)
    private Authority authority;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id")
    private User user;

    @Override
    public String getAuthority() {
        return authority.name();
    }

    @Override
    public final boolean equals(Object o) {
        if (this == o) return true;
        if (o == null) return false;
        Class<?> oEffectiveClass = o instanceof HibernateProxy proxy ?
                proxy.getHibernateLazyInitializer().getPersistentClass() : o.getClass();
        Class<?> thisEffectiveClass = this instanceof HibernateProxy proxy ?
                proxy.getHibernateLazyInitializer().getPersistentClass() : this.getClass();
        if (thisEffectiveClass != oEffectiveClass) return false;
        UserAuthority that = (UserAuthority) o;
        return getId() != null && Objects.equals(getId(), that.getId());
    }

    @Override
    public final int hashCode() {
        return this instanceof HibernateProxy proxy ?
                proxy.getHibernateLazyInitializer().getPersistentClass().hashCode() : getClass().hashCode();
    }

    public UserAuthority(Authority authority, User user) {
        this.authority = authority;
        this.user = user;
    }

    public enum Authority {
        USER,
        ADMIN
    }
}
