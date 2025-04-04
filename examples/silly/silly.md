# Example: Silly

```
schema silly
```
The schema is identified as `silly`. This will be used as the schema name in SQL code

```
enum silly_type {
    "Funny" = 1;
    "Strange" = 2;
    "Dangerous" = 3;
}
```
`silly_type` is an enumeration. It will be implemented as a populated table in [gen.sql](examples/silly/sql/gen.sql).

```

actor {
    id serial PK
    name text not null
}
```
`actor`, `movie` and `movie_actor` specify entities which are generated as tables in 
[gen.sql](examples/silly/sql/gen.sql).
```

movie {
    id serial PK
    name text not null
    silly int FK silly_type.id not null
}
```
movie has many to 1 relation with `silly_type` through the field, `silly` which is a foreign key into `silly_type`.
Field `silly` is mandatory.
```

movie_actor {
    id serial PK
    actor int FK actor.id not null
    movie int FK movie.id not null
}
```
`movie_actor` implements the many-many relationship between `actor` and `movie`.