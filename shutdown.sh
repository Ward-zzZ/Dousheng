#!/bin/bash

# 如果没有参数，则关闭所有服务
if [ $# -eq 0 ]; then
  echo "Stopping all services..."
  killall UserService
  killall RelationService
  killall FavoriteService
  killall CommentService
  killall VideoService
  killall MessageService
  killall api
  echo "All services stopped."
  exit 0
fi

# 如果有参数，则根据参数关闭对应服务
case $1 in
  "U")
    echo "Stopping UserService..."
    killall UserService
    echo "UserService stopped."
    ;;
  "R")
    echo "Stopping RelationService..."
    killall RelationService
    echo "RelationService stopped."
    ;;
  "F")
    echo "Stopping FavoriteService..."
    killall FavoriteService
    echo "FavoriteService stopped."
    ;;
  "C")
    echo "Stopping CommentService..."
    killall CommentService
    echo "CommentService stopped."
    ;;
  "V")
    echo "Stopping VideoService..."
    killall VideoService
    echo "VideoService stopped."
    ;;
  "M")
    echo "Stopping MessageService..."
    killall MessageService
    echo "MessageService stopped."
    ;;
  "A")
    echo "Stopping api..."
    killall api
    echo "api stopped."
    ;;
  *)
    echo "Invalid argument. Usage: ./shutdown.sh [U|R|F|C|V|M|A]"
    exit 1
    ;;
esac

exit 0
