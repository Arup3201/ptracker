import { List, Paragraph, Section } from "../components/legal-content";
import LegalLayout from "../layout/legal-layout";

export default function CookiePage() {
  return (
    <LegalLayout title="Cookie Policy" lastUpdated="April 18, 2026">
      <Paragraph>
        ProjectMate uses cookies to improve your experience.
      </Paragraph>
      <Section title="1. What Are Cookies?">
        <Paragraph>
          Cookies are small text files stored on your device when you visit a
          website.
        </Paragraph>
      </Section>
      <Section title="2. How We Use Cookies">
        <Paragraph>We use cookies for:</Paragraph>
        <List
          items={[
            "Authentication (keeping you logged in)",
            "Session management",
            "Basic functionality of the application",
          ]}
        />
        <Paragraph>
          We do <strong>not</strong> use cookies for:
        </Paragraph>
        <List items={["Advertising", "Third-party tracking", "Analytics"]} />
      </Section>
      <Section title="3. Managing Cookies">
        <Paragraph>
          You can control or disable cookies through your browser settings.
          However, disabling cookies may affect the functionality of the
          application.
        </Paragraph>
      </Section>
      <Section title="4. Changes to This Policy">
        <Paragraph>We may update this Cookie Policy as needed.</Paragraph>
      </Section>
      <Section title="5. Contact">
        <Paragraph>
          For questions, contact us at:{" "}
          <a href="mailto:hello@contact.itsdeployedbyme.dpdns.org">
            hello@contact.itsdeployedbyme.dpdns.org
          </a>
        </Paragraph>
      </Section>
    </LegalLayout>
  );
}
