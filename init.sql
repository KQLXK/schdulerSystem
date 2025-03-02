-- 创建数据库
CREATE DATABASE IF NOT EXISTS schedule;
USE schedule;

-- 创建教师表
CREATE TABLE IF NOT EXISTS teachers (
    id VARCHAR(10) PRIMARY KEY,       -- 教师工号
    name VARCHAR(50) NOT NULL,        -- 教师姓名
    gender VARCHAR(10),               -- 性别
    department VARCHAR(50),           -- 所属院系
    is_external BOOLEAN DEFAULT FALSE, -- 是否外聘
    status VARCHAR(20) DEFAULT '启用', -- 状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
    );

-- 创建教室表
CREATE TABLE IF NOT EXISTS classrooms (
    id VARCHAR(20) PRIMARY KEY,       -- 教室编号
    name VARCHAR(100) NOT NULL,       -- 教室名称
    campus VARCHAR(50),               -- 校区
    building VARCHAR(50),             -- 教学楼
    capacity INT,                     -- 容量
    type VARCHAR(50),                 -- 教室类型（普通教室、多媒体教室等）
    status VARCHAR(20) DEFAULT '启用', -- 状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
    );

-- 创建课程表
CREATE TABLE IF NOT EXISTS courses (
    id VARCHAR(20) PRIMARY KEY,       -- 课程编号
    name VARCHAR(100) NOT NULL,       -- 课程名称
    type VARCHAR(20),                 -- 课程类型（理论、实践、实验）
    credit FLOAT,                     -- 学分
    department VARCHAR(50),           -- 开课院系
    total_hours INT,                  -- 总学时
    status VARCHAR(20) DEFAULT '启用', -- 状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
    );

-- 创建班级表
CREATE TABLE IF NOT EXISTS classes (
    id VARCHAR(20) PRIMARY KEY,       -- 班级编号
    name VARCHAR(100) NOT NULL,       -- 班级名称
    department VARCHAR(50),           -- 所属院系
    major VARCHAR(50),                -- 专业
    campus VARCHAR(50),               -- 校区
    student_count INT,                -- 班级人数
    status VARCHAR(20) DEFAULT '启用', -- 状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
    );

-- 创建排课结果表
CREATE TABLE IF NOT EXISTS schedules (
    id INT AUTO_INCREMENT PRIMARY KEY, -- 排课ID
    course_id VARCHAR(20) NOT NULL,    -- 课程ID
    teacher_id VARCHAR(10) NOT NULL,   -- 教师ID
    classroom_id VARCHAR(20) NOT NULL, -- 教室ID
    class_id VARCHAR(20) NOT NULL,     -- 班级ID
    start_time DATETIME NOT NULL,      -- 开始时间
    end_time DATETIME NOT NULL,        -- 结束时间
    week_pattern VARCHAR(20),          -- 周次模式（单周、双周、全周）
    status VARCHAR(20) DEFAULT '未排',  -- 状态
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    FOREIGN KEY (course_id) REFERENCES courses(id),
    FOREIGN KEY (teacher_id) REFERENCES teachers(id),
    FOREIGN KEY (classroom_id) REFERENCES classrooms(id),
    FOREIGN KEY (class_id) REFERENCES classes(id)
    );