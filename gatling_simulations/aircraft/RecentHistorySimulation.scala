package aircraft

import scala.concurrent.duration._
import scala.util.Random
import java.util.Locale


import io.gatling.core.Predef._
import io.gatling.http.Predef._

class RecentHistorySimulation extends Simulation {
  val BASE_URL = System.getProperty("base_url")
  val N_USERS = Integer.getInteger("users", 20)
  val RAMP = java.lang.Long.getLong("ramp", 0)

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val rand = scala.util.Random
  val resolutionValues = Seq("1h", "6h", "24h")

  val scn = scenario("RecentHistoryScenario")
    .exec{ session =>
        session.set("resolution", resolutionValues(rand.nextInt(resolutionValues.length)))
    }
    .exec(
      http("Airspace recent history request")
        .get(session => ("/airspace/history/realtime/" + session("resolution").as[String]).toLowerCase(java.util.Locale.getDefault))
    )

  setUp(
    scn.inject(
      atOnceUsers(10),
      rampUsers(50).during(10),
      rampUsers(50).during(10),
      constantUsersPerSec(100).during(10))
  ).protocols(httpProtocol)
}
