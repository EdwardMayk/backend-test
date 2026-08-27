FROM nginx:alpine

COPY index.html /usr/share/nginx/html/index.html

EXPOSE 80

CMD ["sh", "-c", "sed -i \"s|__API_BASE_URL__|$API_BASE_URL|g\" /usr/share/nginx/html/index.html && nginx -g \"daemon off;\""]
