# Cotton Cloud Product Requirement Document (PRD)

## 1. Project Overview
**Cotton Cloud (棉花云)** is an AI-powered smart wardrobe management application designed to help users digitize their closet, track their outfit history (OOTD), and make better fashion decisions through intelligent recommendations and data analysis.

The system consists of a native iOS client built with SwiftUI and a backend service powered by Go (Gin framework), utilizing Gemini AI for image processing and fashion analysis.

## 2. Core Features

### 2.1 Smart Closet Management
**Goal**: Digitize the user's physical wardrobe with minimal effort.
*   **Smart Upload**: Users can take photos or upload images from the gallery.
*   **AI Processing**:
    *   **Background Removal**: Automatically removes background from clothing images to create clean, sticker-like assets.
    *   **Auto-Tagging**: AI analyzes the image to suggest Category, Color, Material, Style, and Season tags.
*   **Inventory Tracking**:
    *   **Wear Count**: Tracks how many times an item has been worn.
    *   **Cleanliness Status**: Visual indicators for "Needs Care" based on usage (Max Wear Count vs. Current Wear Count).
    *   **Washing Records**: Log when items are washed.
*   **Organization**:
    *   Filter items by category, color, season, or usage status.
    *   Search/Sort functionality.
    *   Quick edit and bulk delete support.

### 2.2 Outfit of the Day (OOTD) Journal
**Goal**: Track daily style and style evolution.
*   **Calendar View**: Visual calendar showing recorded outfits for past dates.
*   **Outfit Creation**:
    *   Select multiple items from the closet.
    *   **Auto-Collage**: AI generates an aesthetic collage of the selected items.
*   **Contextual Data**:
    *   Automatically records Date.
    *   Integrates Weather data (e.g., "Shanghai · Sunny 22°C") to correlate style with weather conditions.
*   **Stats & Insights**:
    *   "Most Worn" items.
    *   "Total Value" of the wardrobe.
    *   Usage frequency analysis.

### 2.3 Virtual Avatar & Try-On (Experimental)
**Goal**: visualize fit and style.
*   **Avatar Profile**: Users can input body measurements (Height, Weight, Bust, Waist, Hips, etc.).
*   **Virtual Try-On**: AI-driven feature to visualize clothing items on the user's digital avatar.

## 3. User Flows

### 3.1 Adding a New Item (The "Magic" Flow)
1.  **Entry**: User taps "+" in Closet View.
2.  **Capture**: User takes a photo or selects from the gallery.
3.  **Workbench (The Atelier)**:
    *   User sees the raw image.
    *   **Action**: Taps "Magic" button.
    *   **System**: Backend removes background (Cutout) and returns the clean image.
    *   **Refinement** (Optional): User provides text feedback to refine the cutout.
4.  **Classification**:
    *   User reviews/edits AI-suggested tags (Category, Color, Season).
5.  **Save**: Item is saved to the cloud and appears in the Closet grid.

### 3.2 Logging an Outfit
1.  **Entry**: User taps "OOTD" or "Book" icon.
2.  **Date Selection**: Defaults to today or user selects a date.
3.  **Item Selection**: User picks items from the closet.
4.  **Composition**: System generates a collage.
5.  **Save**: Record is saved to the timeline.

## 4. Technical Architecture

### 4.1 Client Side (iOS)
*   **Language**: Swift 5+
*   **Framework**: SwiftUI
*   **Architecture**: MVVM (Model-View-ViewModel)
*   **Key Components**:
    *   `ClosetView`: Main grid, filtering logic.
    *   `ImageWorkbenchView`: Image editing, AI interaction context.
    *   `OOTDCenterView`: Calendar and journal visualization.
    *   `GeminiService`: Handles AI API communication.

### 4.2 Backend Side
*   **Language**: Go (Golang)
*   **Framework**: Gin Web Framework
*   **Database**: SQLite (local dev) / PostgreSQL (production ready) via GORM.
*   **AI Integration**: Google Gemini API for vision tasks.

### 4.3 Data Models

#### User
| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| Nickname | String | Display name |
| Email | String | Account email |

#### ClothingItem
| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| UserID | UUID | Owner ID |
| ImageURL | String | URL of the cropped/clean image |
| OriginalImageURL | String | URL of the raw photo |
| Category | String | e.g., Top, Bottom, Shoes |
| Color | String | Dominant color |
| Season | Array | e.g., ["Summer", "Autumn"] |
| WearCount | Int | Number of times worn |
| MaxWearCount | Int | Hygiene threshold before washing |
| LastWashedAt | Date | Timestamp of last wash |

#### OutfitRecord
| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| UserID | UUID | Owner ID |
| Date | String | "YYYY-MM-DD" |
| Items | Array<UUID> | List of ClothingItem IDs |
| CollageURL | String | Generated visual summary |

#### AvatarProfile
| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| Metrics | JSON | Body measurements (Height, Weight, etc.) |
| ImageURL | String | User's photo or avatar representation |

## 5. API Endpoints

### Authentication
*   `POST /api/v1/auth/register`
*   `POST /api/v1/auth/login`

### Clothing Operations
*   `GET /api/v1/clothing`: List items (filterable).
*   `POST /api/v1/clothing`: Create new item.
*   `PUT /api/v1/clothing/:id`: Update item details.
*   `POST /api/v1/clothing/:id/wear`: Increment wear count.
*   `POST /api/v1/clothing/:id/wash`: Reset wear count / log wash.

### Outfit Operations
*   `GET /api/v1/outfits`: List outfit history.
*   `GET /api/v1/outfits/:date`: Get specific date record.
*   `POST /api/v1/outfits`: Create outfit log.

### AI Services
*   `POST /api/v1/ai/analyze`: Metadata extraction.
*   `POST /api/v1/ai/cutout`: Background removal.
*   `POST /api/v1/ai/refine-cutout`: Interactive refinement.
*   `POST /api/v1/ai/collage`: Outfit composition.
*   `POST /api/v1/ai/tryon`: Virtual try-on generation.
