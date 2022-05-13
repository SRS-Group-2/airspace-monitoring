/*resource "kubernetes_config_map" "flink_config" {
  metadata {
    name = "flink-config"

    labels = {
      app = "flink"
    }
  }

  data = {
    "flink-conf.yaml" = "jobmanager.rpc.address: flink-jobmanager\ntaskmanager.numberOfTaskSlots: 2\nblob.server.port: 6124\njobmanager.rpc.port: 6123\ntaskmanager.rpc.port: 6122\nqueryable-state.proxy.ports: 6125\njobmanager.memory.process.size: 1600m\ntaskmanager.memory.process.size: 1728m\nparallelism.default: 2    \n"

    "log4j-console.properties" = "# This affects logging for both user code and Flink\nrootLogger.level = INFO\nrootLogger.appenderRef.console.ref = ConsoleAppender\nrootLogger.appenderRef.rolling.ref = RollingFileAppender\n\n# Uncomment this if you want to _only_ change Flink's logging\n#logger.flink.name = org.apache.flink\n#logger.flink.level = INFO\n\n# The following lines keep the log level of common libraries/connectors on\n# log level INFO. The root logger does not override this. You have to manually\n# change the log levels here.\nlogger.akka.name = akka\nlogger.akka.level = INFO\nlogger.kafka.name= org.apache.kafka\nlogger.kafka.level = INFO\nlogger.hadoop.name = org.apache.hadoop\nlogger.hadoop.level = INFO\nlogger.zookeeper.name = org.apache.zookeeper\nlogger.zookeeper.level = INFO\n\n# Log all infos to the console\nappender.console.name = ConsoleAppender\nappender.console.type = CONSOLE\nappender.console.layout.type = PatternLayout\nappender.console.layout.pattern = %d{yyyy-MM-dd HH:mm:ss,SSS} %-5p %-60c %x - %m%n\n\n# Log all infos in the given rolling file\nappender.rolling.name = RollingFileAppender\nappender.rolling.type = RollingFile\nappender.rolling.append = false\nappender.rolling.fileName = $${sys:log.file}\nappender.rolling.filePattern = $${sys:log.file}.%i\nappender.rolling.layout.type = PatternLayout\nappender.rolling.layout.pattern = %d{yyyy-MM-dd HH:mm:ss,SSS} %-5p %-60c %x - %m%n\nappender.rolling.policies.type = Policies\nappender.rolling.policies.size.type = SizeBasedTriggeringPolicy\nappender.rolling.policies.size.size=100MB\nappender.rolling.strategy.type = DefaultRolloverStrategy\nappender.rolling.strategy.max = 10\n\n# Suppress the irrelevant (wrong) warnings from the Netty channel handler\nlogger.netty.name = org.jboss.netty.channel.DefaultChannelPipeline\nlogger.netty.level = OFF    \n"
  }
}
*/