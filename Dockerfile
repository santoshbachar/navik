FROM adoptopenjdk/openjdk11:jre-11.0.11_9-alpine

COPY demo.jar /home/demo.jar

expose 8080

#please define java_opts in your environment example JAVA_OPTS=-Xmx512m.

#cmd java $JAVA_OPTS -jar /home/goals.jar
ENTRYPOINT ["sh", "-c",  "exec java $JAVA_OPTS -jar /home/demo.jar"]

