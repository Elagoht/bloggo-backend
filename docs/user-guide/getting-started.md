# Getting Started with Bloggo

## 📋 Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Initial Setup](#initial-setup)
- [User Registration](#user-registration)
- [Creating Your First Post](#creating-your-first-post)
- [Managing Categories and Tags](#managing-categories-and-tags)
- [User Roles and Permissions](#user-roles-and-permissions)
- [Basic Configuration](#basic-configuration)
- [Next Steps](#next-steps)

## 🌟 Introduction

Bloggo is a modern, feature-rich blog management system that makes it easy to create, manage, and publish content. This guide will help you get started with Bloggo and walk you through the essential features.

### Key Features

- **Intuitive Content Management**: Easy-to-use interface for creating and managing blog posts
- **Version Control**: Track changes and manage post versions with approval workflows
- **Rich Media Support**: Upload and manage images, documents, and other media
- **User Management**: Role-based access control for different user types
- **Analytics**: Built-in statistics and analytics to track your blog's performance
- **Responsive Design**: Works seamlessly on desktop and mobile devices

### System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **Memory**: Minimum 4GB RAM, recommended 8GB
- **Storage**: Minimum 2GB free space
- **Network**: Internet connection for initial setup and optional features

## 🚀 Installation

### Option 1: Download Pre-built Binary (Recommended)

1. **Download the Latest Release**
   - Visit the [Bloggo Releases Page](https://github.com/your-org/bloggo/releases)
   - Download the appropriate binary for your operating system
   - Windows: `bloggo-windows.exe`
   - macOS: `bloggo-macos`
   - Linux: `bloggo-linux`

2. **Extract and Prepare**
   ```bash
   # Create a directory for Bloggo
   mkdir bloggo
   cd bloggo

   # Extract the downloaded file
   # (This varies based on your operating system)
   ```

3. **Make the Binary Executable (Linux/macOS)**
   ```bash
   chmod +x bloggo-linux  # or bloggo-macos
   ```

### Option 2: Build from Source

1. **Prerequisites**
   - Install Go 1.24.3 or later from [golang.org](https://golang.org/dl/)
   - Install Git from [git-scm.com](https://git-scm.com/)

2. **Clone and Build**
   ```bash
   git clone https://github.com/your-org/bloggo.git
   cd bloggo
   go build -o bloggo ./cli
   ```

### Option 3: Docker (Advanced)

1. **Using Docker Compose**
   ```bash
   # Create docker-compose.yml
   cat > docker-compose.yml << EOF
   version: '3.8'
   services:
     bloggo:
       image: bloggo:latest
       ports:
         - "8723:8723"
       volumes:
         - ./data:/app/data
   EOF

   # Start Bloggo
   docker-compose up -d
   ```

## ⚙️ Initial Setup

### 1. Start Bloggo

Run the Bloggo binary:

```bash
# Using the downloaded binary
./bloggo

# Or if built from source
go run ./cli/main.go
```

### 2. First-Time Configuration

When you run Bloggo for the first time, it will:

1. **Generate Configuration File**: Creates `bloggo-config.json` with secure random values
2. **Initialize Database**: Creates `bloggo.sqlite` with the required tables
3. **Create Default Admin User**: Sets up an administrator account
4. **Start the Server**: Begins listening on port 8723

You should see output similar to:

```
Starting server on http://localhost:8723
Configuration file created: bloggo-config.json
Database initialized: bloggo.sqlite
Default admin user created
Server is ready to accept connections
```

### 3. Access Bloggo

Open your web browser and navigate to:

```
http://localhost:8723
```

You should see the Bloggo welcome page or login interface.

## 👤 User Registration

### Default Admin Account

Bloggo automatically creates a default admin user:

- **Email**: `admin@bloggo.local`
- **Password**: `admin123` (change this immediately!)

### Creating Your First User Account

1. **Navigate to Registration**
   - Go to `http://localhost:8723/register`
   - Or click "Register" on the login page

2. **Fill Registration Form**
   ```
   Name: John Doe
   Email: john.doe@example.com
   Password: YourSecurePassword123
   Confirm Password: YourSecurePassword123
   ```

3. **Verify Email** (if configured)
   - Check your email for verification link
   - Click the link to activate your account

### Changing Default Admin Password

1. **Login as Admin**
   - Use the default admin credentials

2. **Go to Profile Settings**
   - Click on your username/profile
   - Select "Settings" or "Profile"

3. **Update Password**
   ```
   Current Password: admin123
   New Password: YourNewSecurePassword123
   Confirm Password: YourNewSecurePassword123
   ```

## ✍️ Creating Your First Post

### Step 1: Access the Post Editor

1. **Login to Bloggo**
   - Use your credentials to login

2. **Navigate to Posts**
   - Click "Posts" in the navigation menu
   - Click "New Post" button

### Step 2: Write Your Post

1. **Post Information**
   ```
   Title: My First Blog Post
   Slug: my-first-blog-post  (auto-generated)
   Status: Draft  (will change to Published later)
   ```

2. **Post Content**
   Bloggo supports Markdown for content formatting:

   ```markdown
   # My First Blog Post

   This is the introduction to my first blog post. I'm excited to share my thoughts with the world!

   ## What I'll Cover

   - My background
   - Why I started blogging
   - What you can expect from future posts

   ## My Background

   I've always been passionate about [technology and writing](https://example.com).
   This blog will be my space to explore ideas and share knowledge.

   ## Why I Started Blogging

   Blogging gives me a platform to:
   1. Share my expertise
   2. Connect with others
   3. Document my learning journey

   Thanks for reading!
   ```

3. **Excerpt (Optional)**
   ```
   A brief introduction to my blogging journey and what you can expect from future posts.
   ```

### Step 3: Add Media

1. **Upload Cover Image**
   - Click "Upload Cover Image"
   - Select an image from your computer
   - Bloggo will automatically resize and optimize it

2. **Insert Images in Content**
   ```markdown
   ![Alt text for image](/uploads/your-image.jpg)
   ```

### Step 4: Categorize and Tag

1. **Select Categories**
   - Choose "Technology" (or create a new category)
   - Categories help organize your content

2. **Add Tags**
   - Add relevant tags: `blogging`, `first-post`, `introduction`
   - Tags help with content discovery

### Step 5: Review and Publish

1. **Preview Your Post**
   - Click "Preview" to see how it will look
   - Check formatting and content

2. **Publish**
   - Change status from "Draft" to "Published"
   - Click "Save" or "Publish"

3. **View Your Post**
   - Your post is now live at:
   ```
   http://localhost:8723/posts/my-first-blog-post
   ```

## 🏷️ Managing Categories and Tags

### Creating Categories

1. **Navigate to Categories**
   - Click "Categories" in the admin menu

2. **Create New Category**
   ```
   Name: Technology
   Slug: technology  (auto-generated)
   Description: Posts about technology, programming, and digital innovation
   Parent Category: (None for top-level categories)
   ```

3. **Create Subcategories**
   ```
   Name: Go Programming
   Slug: go-programming
   Description: Posts specifically about Go programming language
   Parent Category: Technology
   ```

### Managing Tags

1. **Navigate to Tags**
   - Click "Tags" in the admin menu

2. **Create New Tag**
   ```
   Name: Go Programming
   Slug: go-programming  (auto-generated)
   ```

3. **Tag Management**
   - View all tags and their usage counts
   - Edit or delete unused tags
   - Merge similar tags

### Best Practices

- **Categories**: Use broad topics, keep the number manageable (5-10 main categories)
- **Tags**: Be specific, use as many as needed, but be consistent
- **Hierarchy**: Use subcategories for better organization
- **Naming**: Use clear, descriptive names

## 👥 User Roles and Permissions

### Understanding User Roles

Bloggo has three main user roles:

#### Administrator (Admin)
- **Full Access**: Can manage all aspects of the blog
- **User Management**: Create, edit, and delete user accounts
- **System Settings**: Configure blog settings and features
- **Content Control**: Full control over all posts and content

#### Editor
- **Content Management**: Can edit, publish, and delete any posts
- **Category/Tag Management**: Create and manage categories and tags
- **User Permissions**: Limited user management capabilities
- **Analytics Access**: Can view blog statistics and analytics

#### Author
- **Own Content**: Can create, edit, and delete their own posts
- **Draft Management**: Can save drafts and submit for approval
- **Media Upload**: Can upload images and files for their posts
- **Limited Access**: Cannot edit others' content or system settings

### User Management

#### Creating New Users

1. **Navigate to Users**
   - Click "Users" in the admin menu

2. **Add New User**
   ```
   Name: Jane Smith
   Email: jane.smith@example.com
   Role: Author
   ```

3. **Set Initial Password**
   - System will generate a temporary password
   - User will be prompted to change on first login

#### Managing User Permissions

1. **Edit User Profile**
   - Click on a user in the user list
   - Modify their information and role

2. **Change User Role**
   - Select appropriate role from dropdown
   - Changes take effect immediately

### Permission Matrix

| Feature | Admin | Editor | Author |
|---------|-------|--------|--------|
| Create Posts | ✅ | ✅ | ✅ |
| Edit Own Posts | ✅ | ✅ | ✅ |
| Edit All Posts | ✅ | ✅ | ❌ |
| Delete Posts | ✅ | ✅ | ✅ (own only) |
| Manage Categories | ✅ | ✅ | ❌ |
| Manage Users | ✅ | ❌ | ❌ |
| View Analytics | ✅ | ✅ | ✅ (limited) |
| System Settings | ✅ | ❌ | ❌ |

## ⚙️ Basic Configuration

### General Settings

1. **Access Settings**
   - Click "Settings" in the admin menu
   - Navigate to "General" tab

2. **Basic Configuration**
   ```
   Blog Title: My Awesome Blog
   Blog Description: A blog about technology and life
   Admin Email: admin@example.com
   Timezone: UTC
   Date Format: January 2, 2006
   ```

### Appearance Settings

1. **Theme Selection**
   - Navigate to "Appearance" tab
   - Choose from available themes
   - Customize colors and fonts

2. **Logo and Branding**
   - Upload blog logo
   - Set favicon
   - Configure header and footer

### Comment Settings

1. **Enable Comments**
   ```
   Enable Comments: Yes
   Comment Moderation: Required for new users
   Comment Expiration: Never
   ```

2. **Comment Form Fields**
   ```
   Name: Required
   Email: Required
   Website: Optional
   Comment: Required
   ```

### SEO Settings

1. **Basic SEO**
   ```
   Meta Description: A blog about technology and programming
   Meta Keywords: technology, programming, blog
   Google Analytics: (tracking code)
   ```

2. **Social Media**
   ```
   Facebook Page URL: https://facebook.com/yourblog
   Twitter Handle: @yourblog
   LinkedIn Company Page: https://linkedin.com/company/yourblog
   ```

## 🔧 Advanced Configuration

### Security Settings

1. **Password Policy**
   ```
   Minimum Password Length: 8
   Require Uppercase: Yes
   Require Lowercase: Yes
   Require Numbers: Yes
   Require Special Characters: No
   ```

2. **Session Management**
   ```
   Session Timeout: 24 hours
   Maximum Concurrent Sessions: 3
   Force Logout on Password Change: Yes
   ```

### Email Configuration

1. **SMTP Settings**
   ```
   SMTP Host: smtp.gmail.com
   SMTP Port: 587
   SMTP Username: your-email@gmail.com
   SMTP Password: your-app-password
   Use TLS: Yes
   ```

2. **Email Templates**
   - Configure welcome emails
   - Set up password reset templates
   - Customize notification emails

### Backup Configuration

1. **Automatic Backups**
   ```
   Enable Backups: Yes
   Backup Frequency: Daily
   Backup Retention: 30 days
   Backup Location: ./backups/
   ```

2. **Manual Backup**
   - Click "Backup Now" button
   - Download backup file
   - Schedule regular backups

## 📊 Monitoring Your Blog

### Dashboard Overview

1. **Access Dashboard**
   - Click "Dashboard" in admin menu
   - View key metrics at a glance

2. **Key Metrics**
   - Total posts: 25
   - Published posts: 20
   - Total views: 1,250
   - Users: 15
   - Comments: 45

### Analytics and Statistics

1. **Post Performance**
   - Most viewed posts
   - Recent views
   - Engagement metrics

2. **User Analytics**
   - New registrations
   - Active users
   - Login statistics

3. **Traffic Sources**
   - Direct traffic
   - Referral sources
   - Search engine traffic

## 🚀 Next Steps

### Content Strategy

1. **Plan Your Content Calendar**
   - Decide on posting frequency (e.g., weekly)
   - Create a content calendar
   - Plan post topics in advance

2. **Develop Your Voice**
   - Define your target audience
   - Establish your writing style
   - Create content guidelines

3. **Engage with Your Audience**
   - Respond to comments
   - Encourage discussion
   - Build a community

### Technical Improvements

1. **Customize Your Blog**
   - Install themes and plugins
   - Customize the appearance
   - Add custom features

2. **Optimize Performance**
   - Enable caching
   - Optimize images
   - Monitor page load times

3. **SEO Optimization**
   - Research keywords
   - Optimize meta tags
   - Build backlinks

### Security and Maintenance

1. **Regular Updates**
   - Keep Bloggo updated
   - Update plugins and themes
   - Monitor security advisories

2. **Backup Strategy**
   - Regular backups
   - Test restore procedures
   - Off-site backup storage

3. **Security Best Practices**
   - Use strong passwords
   - Enable two-factor authentication
   - Monitor user activity

### Growing Your Blog

1. **Promote Your Content**
   - Share on social media
   - Engage with other bloggers
   - Participate in online communities

2. **Build Your Audience**
   - Create an email newsletter
   - Engage with readers
   - Collaborate with other creators

3. **Monetization Options**
   - Display advertising
   - Affiliate marketing
   - Sponsored content
   - Digital products

## 🆘 Getting Help

### Documentation

- **User Guide**: This document and related guides
- **API Documentation**: For developers and integrations
- **FAQ**: Frequently asked questions

### Community Support

- **GitHub Issues**: Report bugs and request features
- **Discussion Forums**: Connect with other users
- **Stack Overflow**: Technical questions

### Professional Support

- **Documentation**: Comprehensive guides and tutorials
- **Community Forums**: Get help from other users
- **Professional Services**: Paid support and consulting

### Contact Information

- **Email**: support@bloggo.dev
- **Website**: https://bloggo.dev
- **GitHub**: https://github.com/your-org/bloggo

---

**Document Version**: 1.0.0
**Last Updated**: October 4, 2025
**Author**: Bloggo Documentation Team
**Reviewers**: User Experience Committee