## PRD (Product Requirements Document)

## future features:

- Implement a search functionality to allow users to search for listings based on keywords, categories, and locations.
- Implement a rating and review system for listings, allowing users to provide feedback and ratings for products or services.
- Implement a wish list feature, allowing users to save and manage their favorite listings.

## other features:

- Only admin user can create categories.

```sql
ALTER TABLE users
ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

ALTER TABLE users
ADD CONSTRAINT users_role_check
CHECK (role IN ('user', 'admin'));
```
