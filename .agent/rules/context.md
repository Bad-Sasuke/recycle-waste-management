---
trigger: always_on
---

# Project Guidelines & Architecture

## 📋 Best Practices

### Frontend Development Guidelines

Before creating or modifying any page, ensure the following:

- Follow the coding standards outlined in [Clean Code TypeScript](https://github.com/labs42io/clean-code-typescript)
- Keep code comments minimal — only add them when truly necessary

#### Checklist:

1. **Design Tokens** — Use existing tokens from `frontend/src/assets/main.css` (based on TailwindCSS)
2. **Components** — Check `src/components/` for reusable components → [Vue Best Practices](https://learnvue.co/articles/vue-best-practices)
- ทุกครั้งที่สร้าง component ใหม่ ให้สร้าง story book ด้วย
- ทุกครั้งที่มีการแก้ไข component ให้แก้ story book ด้วย
- How to create a story book so that the paths must be in the same folder.
3. **State Management** — Use existing stores in `frontend/src/stores` → [Pinia Best Practices](https://masteringpinia.com/blog/5-best-practices-for-scalable-vuejs-state-management-with-pinia)
4. **API Integration** — Follow the existing structure: `types/` (response/payload types) + `services/` (API calls)
5 check and fixed wih command `npm run type-check` then `npm run lint` everytime 
---

## 🎨 Frontend Features

### Design System (`frontend/src/assets/`)

- Design tokens and CSS variables (TailwindCSS-based)
- Logo and static assets

### UI Components (`frontend/src/components/`)

- **Navigation & Layout**: Header navigation, Footer, Glassmorphism cards, Navigation buttons
- **Authentication & User**: Login popup, Role selection modal, Profile management
- **Marketplace**: Product catalog, Product detail modal, Landing page
- **Waste Management**: Waste item popups, Edit waste forms, Category cards
- **Map & Location** (`components/map/`): Customer markers, Shop markers, User location markers, Animated liquid markers
- **Notifications**: Shop notifications, Alert modals, Error messages
- **i18n** (`components/switchLang/`): Language switcher (Thai/English)

### Pages (`frontend/src/views/`)

- **Public**: Home, About, Marketplace
- **Authentication**: OAuth login/callback
- **User Management**: Profile, Settings
- **Shop Management**:
  - Create shop (moderator)
  - Manage shop dashboard (inventory, pricing, reviews)
  - Shop profile (public view)
- **Location & Map**:
  - Shop locator (find nearby shops)
  - Share location (customer request pickup)
  - Nearby customers (shop owner view requests)
- **Transaction**: Waste purchase/payment processing

### State Management (`frontend/src/stores/`)

- User authentication & session
- Shop data & inventory
- Recyclable waste items & categories
- Customer pickup requests
- Purchase transactions
- Real-time chat (WebSocket)
- Cookie/session management
- i18n language state

### API Services (`frontend/src/services/`)

- Base HTTP client configuration
- Authentication (OAuth GitHub/Google)
- Shop CRUD operations
- Shop management (inventory, pricing)
- Customer request handling
- Firebase integration

### Type Definitions (`frontend/src/types/`)

- API response/request types
- Shop, waste, category, customer request entities
- Type-safe API integration

### Routing (`frontend/src/router/`)

- Route definitions & navigation
- Authentication guards
- Role-based access control

### Internationalization (`frontend/src/i18n/`)

- Thai/English language support
- Translation files

---

## 🔧 Backend Features

### Development Standards

- Follow the coding patterns from [Clean Go Article](https://github.com/Pungyeon/clean-go-article)
- Keep code comments minimal — only add them when truly necessary
- update swagger.yaml ให้ตรงกับ api หากมีการแก้ไขหรือเพิ่ม api ใหม่
- check and fixed error wih command `cd backend && go build main.go`

### Configuration (`backend/src/configuration/`)

- Fiber web framework setup
- Environment configuration

### Domain Layer (`backend/src/domain/`)

- Data Transfer Objects (`domain/entities/`)
- Business models (`domain/models/`)
- Redis cache interfaces (`domain/caches/`)
- Database interfaces (`domain/datasources/`)

### API Endpoints (`backend/src/gateways/`)

**Services**: User, Recycle Waste, Auth, Image, Shop, Settings, Customer Request

**Available APIs**:

- **User Management**: Create account, profile CRUD, role management, profile image upload
- **Waste Management**: List recyclable items, add/edit/delete waste items
- **Categories**: Get waste categories
- **Authentication**: OAuth callbacks (GitHub, Google), JWT token generation
- **Shop**: Create/manage shops, public shop listing, shop profiles
- **Settings**: User preferences management
- **Customer Requests**: Create/manage pickup requests, accept/cancel/complete requests, walk-in customers
- **Real-time Chat**: WebSocket connection for shop-customer communication
- **Reviews**: Submit reviews, skip review, check review status
- **Receipts**: Generate receipts, view shop receipts
- **Stock Management**: View shop inventory

### Infrastructure (`backend/src/infrastructure/`)

- External HTTP clients (`infrastructure/httpclients/`)
- Third-party service providers (`infrastructure/providers/`)
- Utility functions (`infrastructure/utils/`)

### Middleware (`backend/src/middlewares/`)

- JWT authentication verification
- Request/response logging
- Performance monitoring
- API documentation (Scalar)

### Data Access (`backend/src/repositories/`)

- Database CRUD operations
- Query builders
- Data persistence

### Business Logic (`backend/src/services/`)

- Authentication logic
- User management
- Shop operations
- Customer request processing
- Receipt generation
- Waste recycling logic
- Review management
- Settings management
- Stock/inventory management
- Image processing

### Real-time Communication (`backend/src/websocket/`)

- WebSocket hub (connection management)
- Message broadcasting
- Online/offline status
- Chat message handling

### Documentation (`backend/src/docs/`)

- API specifications
- OpenAPI/Swagger docs

---