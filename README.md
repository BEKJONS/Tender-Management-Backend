# Tender Management Backend

## Project Overview
The Tender Management Backend system enables clients to post tenders, contractors to submit bids, and clients to review and award contracts. It includes features such as user authentication, role management, tender creation, bid submission, bid evaluation, and real-time notifications.

## Features
- **User Authentication & Role Management**: Registration, login, and role-based access.
- **Tender Posting**: Clients can create, list, and manage tenders.
- **Bid Submission**: Contractors can submit bids with price, delivery time, and comments.
- **Bid Evaluation & Tender Awarding**: Clients evaluate bids and award tenders.
- **Real-time Notifications**: WebSockets for real-time updates.

## Database Schema
- **User**: `id, username, password, role, email`
- **Tender**: `id, client_id, title, description, deadline, budget, status`
- **Bid**: `id, tender_id, contractor_id, price, delivery_time, comments, status`
- **Notification**: `id, user_id, message, relation_id, type, created_at`

## Setup Instructions

### Database Setup
Start the database with Docker:
```bash
make run_db
