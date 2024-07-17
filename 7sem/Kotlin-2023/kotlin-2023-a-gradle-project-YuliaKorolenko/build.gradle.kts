plugins {
    alias(libs.plugins.kotlin)
    alias(libs.plugins.ktlintG)
    alias(libs.plugins.detekt)
}

group = "org.example"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    testImplementation(libs.kotest.runner)
    testImplementation(libs.mockk)
    implementation(libs.ktlint)
}

tasks.withType<Test>().configureEach {
    useJUnitPlatform()
}

detekt {
    config.setFrom("$projectDir/detekt.yml")
}

ktlint {
    disabledRules.set(setOf("no-wildcard-imports"))
}
