# 🌸✨ K-pop Photocard Collection ✨🌸

A beautiful, girly web application for managing your precious K-pop photocard collection! Built with Go backend and a cute pink aesthetic that'll make your heart flutter! 💕

![K-pop PC Collection](https://img.shields.io/badge/K--pop-Photocard%20Collection-ff69b4?style=for-the-badge&logo=music&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white)

## 🎀 Features

### ✨ Core Functionality
- **📸 Photo Upload**: Upload photocard images with drag & drop support
- **🏷️ Auto-Tagging**: Automatic hashtag generation based on group, album, and member
- **🔍 Advanced Search**: Search across groups, members, albums, and tags
- **🎛️ Filter Sidebar**: Filter by groups, members, and albums with checkboxes
- **📝 Edit & Delete**: Click any card to edit details or remove from collection
- **🎨 Realistic Display**: Photocards displayed in accurate 55x85mm proportions
- **📱 Responsive Design**: Works beautifully on desktop, tablet, and mobile

### 💖 User Experience
- **Girly Pop Aesthetic**: Pink gradients, cute shadows, and sparkle animations
- **Auto-Complete**: Smart group name suggestions while typing
- **Real-time Filtering**: Instant search results as you type
- **Hover Effects**: Cute animations when interacting with cards
- **Loading States**: Beautiful loading indicators and error messages

### 🔧 Technical Features
- **Lightweight Backend**: Pure Go with no external dependencies
- **JSONL Database**: Simple file-based storage that's easy to backup
- **RESTful API**: Clean API endpoints for frontend interactions
- **File Management**: Automatic image file handling and organization
- **Error Handling**: Comprehensive error handling and user feedback

## 📁 Project Structure

```
kpop-pc-site/
├── main.go                 # 🚀 Main Go server (heavily commented!)
├── cards.jsonl            # 💾 Database file (JSON Lines format)
├── templates/
│   └── index.html         # 🎨 Frontend template (with extensive JS comments)
├── static/                # 📸 Uploaded photocard images stored here
│   ├── 1759803801_Felix_StrayKids.jpg
│   ├── 1759804640_Chaeyoung_Twice.jpg
│   └── ...more images...
└── README.md              # 📖 This documentation
```

## 🚀 Quick Start

### Prerequisites
- **Go 1.16 or higher** ([Download Go](https://golang.org/dl/))
- **A modern web browser** (Chrome, Firefox, Safari, Edge)
- **Some cute photocard images** to upload! 📸

### Installation

1. **Clone or download this repository**
   ```bash
   git clone <repository-url>
   cd kpop-pc-site
   ```

2. **Create necessary directories**
   ```bash
   mkdir -p static templates
   ```

3. **Start the server**
   ```bash
   go run main.go
   ```

4. **Open your browser and visit**
   ```
   http://localhost:8080
   ```

5. **Start collecting! 🎉**

### First Time Setup

The application will automatically create:
- `cards.jsonl` - Your photocard database
- `static/` - Directory for uploaded images

## 📖 How to Use

### 🆕 Adding Your First Photocard

1. **Click "Add New Photocard ✨"** on the main page
2. **Fill out the form:**
   - **Group**: Start typing and see auto-complete suggestions! (e.g., "ATEEZ", "NewJeans")
   - **Album**: Album name (e.g., "Zero : Fever Pt.2", "Get Up")
   - **Member**: Member name (e.g., "Mingi", "Hanni")
   - **Copies**: How many copies you own (default: 1)
   - **Photo**: Upload your photocard image (JPG, PNG, etc.)
3. **Click "Upload Photocard ✨"**
4. **Watch your new card appear in the collection!** 🎊

### 🔍 Finding Your Cards

#### Search Bar
- **Type anything** in the search bar to find cards instantly
- **Searches across**: Group names, member names, album names, and tags
- **Example searches**: "ATEEZ", "Fever", "Mingi", "#NewJeans"

#### Filter Sidebar
- **Group Filter**: Check boxes to show only specific groups
- **Member Filter**: Filter by your favorite members
- **Album Filter**: See cards from specific albums
- **Combine Filters**: Use multiple filters together!
- **Clear All**: Reset all filters with one click

### ✏️ Editing Cards

1. **Click on any photocard** in your collection
2. **Edit the details** in the popup modal
3. **Click "Save Changes ✨"** to update
4. **Or click "Delete Forever 💔"** to remove (with confirmation!)

## 🗂️ File Storage System

### 📸 Image Storage
- **Location**: All uploaded images are stored in the `static/` directory
- **Naming**: Files are renamed with timestamp prefixes to prevent conflicts
  - Format: `{timestamp}_{original_filename}`
  - Example: `1759803801_Felix_StrayKids.jpg`
- **Access**: Images are served at `/static/{filename}` URLs
- **Formats**: Supports JPG, PNG, GIF, and other web-compatible formats

### 💾 Database Storage
- **File**: `cards.jsonl` (JSON Lines format)
- **Format**: One JSON object per line for easy parsing and appending
- **Example entry**:
  ```json
  {\"group\":\"ATEEZ\",\"album\":\"Zero : Fever Pt.2\",\"member\":\"Mingi\",\"copies\":2,\"image\":\"/static/1759803801_mingi.jpg\",\"tags\":[\"#ATEEZ\",\"#ZeroFeverPt2\",\"#Mingi\"]}
  ```
- **Backup**: Simply copy the `cards.jsonl` file to backup your collection data
- **Restore**: Replace `cards.jsonl` with your backup to restore

## 🎨 Customization

### 🌈 Changing Colors
The girly pink theme can be customized by editing the CSS in `templates/index.html`:

```css
.girly-gradient {
  background: linear-gradient(135deg, #ff9a9e 0%, #fecfef  100%);
}
```

### 📐 Card Proportions
Photocard aspect ratio is set to real-world proportions (55x85mm):
```css
.aspect-[55/85]  /* Width 55, Height 85 */
```

### 🔤 Fonts and Icons
- **CSS Framework**: TailwindCSS 2.2.19
- **Icons**: Font Awesome 6.0.0
- **Fonts**: System fonts for best performance

## 🔌 API Reference

The backend provides RESTful endpoints for the frontend:

### GET `/api/cards`
Returns all photocard data as JSON array
```json
[
  {
    \"group\": \"ATEEZ\",
    \"album\": \"Zero : Fever Pt.2\", 
    \"member\": \"Mingi\",
    \"copies\": 2,
    \"image\": \"/static/1759803801_mingi.jpg\",
    \"tags\": [\"#ATEEZ\", \"#ZeroFeverPt2\", \"#Mingi\"]
  }
]
```

### GET `/api/groups`
Returns unique group names for auto-complete
```json
[\"ATEEZ\", \"NewJeans\", \"LE SSERAFIM\", \"Stray Kids\"]
```

### POST `/api/update`
Update an existing photocard
```json
{
  \"index\": 0,
  \"card\": {
    \"group\": \"ATEEZ\",
    \"album\": \"Updated Album\",
    \"member\": \"Mingi\",
    \"copies\": 3,
    \"image\": \"/static/existing_image.jpg\",
    \"tags\": [\"#ATEEZ\", \"#UpdatedAlbum\", \"#Mingi\"]
  }
}
```

### POST `/api/delete`
Delete a photocard by index
```json
{\"index\": 0}
```

## 🚀 Deployment Options

### 🏠 Local Development
```bash
go run main.go
# Server runs on http://localhost:8080
```

### 🌐 Production Deployment

#### Option 1: Build Binary
```bash
# Build for your platform
go build -o photocard-server main.go

# Run the binary
./photocard-server
```

#### Option 2: Docker (create your own Dockerfile)
```dockerfile
FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD [\"./main\"]
```

#### Option 3: Cloud Platforms
- **Heroku**: Add a `Procfile` with `web: ./main`
- **Railway**: Works out of the box with Go detection
- **Digital Ocean**: Deploy as a simple droplet

### 🔒 Security Considerations
- **File Upload Limits**: Currently set to 10MB per image
- **CORS**: Not configured - add if serving from different domains
- **HTTPS**: Not implemented - add TLS certificates for production
- **Authentication**: No user system - anyone with access can modify collection

## 🐛 Troubleshooting

### Common Issues

#### ❌ \"Server won't start\"
```bash
# Check if port 8080 is already in use
lsof -i :8080

# Kill existing process if needed
kill -9 <PID>
```

#### ❌ \"Images not displaying\"
- Check that `static/` directory exists and has proper permissions
- Verify image file extensions are supported
- Check browser console for 404 errors

#### ❌ \"Filters showing 'undefined'\"
- Open browser developer tools and check console for JavaScript errors
- Verify `/api/cards` endpoint returns valid JSON
- Check network tab for failed API requests

#### ❌ \"Upload form not working\"
- Ensure you're using POST method
- Check file size is under 10MB limit
- Verify all required fields are filled

### 🔍 Debug Mode
Add debug logging to `main.go`:
```go
log.Printf(\"Debug: %+v\", yourVariable)
```

## 🤝 Contributing

Want to make this even more cute? Here are some ideas:

### 🌟 Feature Ideas
- **🔐 User Authentication**: Multiple collections per user
- **📊 Statistics**: Collection stats and insights
- **🎵 Audio Support**: Add voice clips or song previews
- **📱 PWA Support**: Make it installable as a mobile app
- **🌙 Dark Mode**: For late-night collecting sessions
- **📤 Export/Import**: Backup collections to different formats
- **🔗 Social Features**: Share collections with friends

### 🎨 UI Improvements  
- **🌈 Theme Options**: Multiple color schemes
- **✨ More Animations**: Even more delightful interactions
- **📐 Grid Options**: Different layout sizes
- **🖼️ Image Editor**: Crop and filter images before upload
- **💫 Transitions**: Smoother page transitions

### 🏗️ Technical Enhancements
- **🗄️ Database Options**: SQLite or PostgreSQL support
- **🔄 Real-time Updates**: WebSocket support for live updates
- **🚀 Performance**: Image optimization and caching
- **📱 API Improvements**: Pagination and sorting
- **🧪 Testing**: Unit and integration tests

## 📄 License

This project is open source and available under the MIT License. Feel free to use it for your own photocard collections! 💕

## 💕 Acknowledgments

- **K-pop Community**: For inspiring this project with amazing photocard collections
- **Go Community**: For creating such a wonderful language
- **TailwindCSS**: For making beautiful styling so easy
- **Font Awesome**: For the cute icons

---

## 🌸 Final Notes

This project was created with love for the K-pop community and photocard collectors everywhere! 

Whether you're collecting your first card or have thousands of precious memories, this app will help you organize and enjoy your collection in the most adorable way possible! ✨

**Happy Collecting! 🎊💖**

---

*Made with 💕 for K-pop fans worldwide*