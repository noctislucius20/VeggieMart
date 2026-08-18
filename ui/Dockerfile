FROM oven/bun:1-alpine

WORKDIR /app

COPY package*.json ./
COPY bun.lock ./

RUN bun install

COPY . .

ENV PORT=3000
ENV NODE_ENV=production

RUN bun run build

EXPOSE 3000

CMD ["bun", ".output/server/index.mjs"]