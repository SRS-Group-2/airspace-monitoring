package aircraft

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._

class InfoSimulation extends Simulation {
  val BASE_URL = "https://aircraft-info-74anfdzq5q-uc.a.run.app"
  val N_USERS = 600

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val scn = scenario("InfoSimulation")
    .exec(
      http("123456 aircraft info request")
        .get("/airspace/aircraft/123456/info")
    )

  setUp(
    scn.inject(atOnceUsers(N_USERS))
  ).protocols(httpProtocol)
}
