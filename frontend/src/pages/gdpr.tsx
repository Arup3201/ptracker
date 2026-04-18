import { List, Paragraph, Section } from "../components/legal-content";
import LegalLayout from "../layout/legal-layout";

export default function GDPRPage() {
  return (
    <LegalLayout title="GDPR Notice" lastUpdated="April 18, 2026">
      <Paragraph>
        If you are located in the European Economic Area (EEA), you have certain
        data protection rights under the General Data Protection Regulation
        (GDPR).
      </Paragraph>
      <Section title="Your Rights">
        <Paragraph>You have the right to:</Paragraph>
        <List
          items={[
            "Access your personal data",
            "Correct inaccurate data",
            "Request deletion of your data",
            "Restrict or object to processing",
            "Request data portability",
          ]}
        />
      </Section>
      <Section title="Legal Basis for Processing">
        <Paragraph>We process your data based on:</Paragraph>
        <List
          items={[
            "Your consent (when signing up)",
            "Necessity to provide the service",
          ]}
        />
      </Section>
      <Section title="Data Transfers">
        <Paragraph>
          Your data may be stored and processed outside the EEA (e.g., in
          India). We take reasonable steps to ensure your data is handled
          securely.
        </Paragraph>
      </Section>
      <Section title="Contact">
        <Paragraph>
          To exercise your rights, contact:{" "}
          <a href="mailto:hello@contact.itsdeployedbyme.dpdns.org">
            hello@contact.itsdeployedbyme.dpdns.org
          </a>
        </Paragraph>
      </Section>
    </LegalLayout>
  );
}
